package controllers

import (
	"fmt"
	"gitlab.com/aoterocom/changelog-guardian/application/helpers"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	services2 "gitlab.com/aoterocom/changelog-guardian/application/services"
	"gitlab.com/aoterocom/changelog-guardian/controller/interfaces"
	"gitlab.com/aoterocom/changelog-guardian/controller/services"
	infraInterfaces "gitlab.com/aoterocom/changelog-guardian/infrastructure/interfaces"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"strconv"
	"time"
)

type ChangelogGuardianController struct {
	releaseProvider infraInterfaces.Provider
	taskProvider    infraInterfaces.Provider
	releasePipes    []interfaces.ReleasePipe
	taskPipes       []interfaces.TaskPipe
}

func NewChangelogGuardianController(releaseProvider infraInterfaces.Provider, taskProvider infraInterfaces.Provider, releasePipesStr []string, taskPipesStr []string) (*ChangelogGuardianController, error) {
	// reversing task and release pipes to start iterating over the last
	helpers.ReverseAny(releasePipesStr)
	helpers.ReverseAny(taskPipesStr)

	var releasePipes []interfaces.ReleasePipe
	for _, pipeStr := range releasePipesStr {
		pipe, err := services.ReleasePipeSelector(pipeStr)
		if err != nil {
			return nil, err
		}
		releasePipes = append(releasePipes, *pipe)
	}

	var taskPipes []interfaces.TaskPipe
	for _, pipeStr := range taskPipesStr {
		pipe, err := services.TaskPipeSelector(pipeStr)
		if err != nil {
			return nil, err
		}
		taskPipes = append(taskPipes, *pipe)
	}

	return &ChangelogGuardianController{releaseProvider: releaseProvider, taskProvider: taskProvider, releasePipes: releasePipes, taskPipes: taskPipes}, nil
}

func (cgc *ChangelogGuardianController) GetFilledReleasesFromInfra(lastRelease *models.Release, mainBranch string, defaultBranch string) (*[]models.Release, error) {
	var from1 *time.Time
	if lastRelease != nil {
		layout := "2006-01-02T15:04:05"
		str := lastRelease.Date + "T00:00:00"
		t, _ := time.Parse(layout, str)
		from1 = &t
	} else {
		from1 = nil
	}

	releases, err := cgc.releaseProvider.GetReleases(from1, nil)
	if err != nil {
		return nil, err
	}

	var infraTruncatedReleases []infra.Release
	var appTruncatedReleases []models.Release

	// If there is a release from where to search, we truncate the releases obtained just from it
	if lastRelease != nil {
		for i, release := range *releases {
			if release.Name == lastRelease.Version {
				helpers.ReverseAny(*releases)
				if len((*releases)[:i+1]) == 1 {
					infraTruncatedReleases = []infra.Release{}
				} else {
					infraTruncatedReleases = (*releases)[:i+1]
				}
				break
			}
		}
		// If not, just use all the releases
	} else {
		infraTruncatedReleases = *releases
	}

	// Pass releases through Release Pipes
	infraTruncatedReleases = cgc.throughReleasePipes(infraTruncatedReleases)

	// For each release obtained from infra layer
	for i, release := range infraTruncatedReleases {

		// Map the release to application layer model
		appTruncatedRelease := *services.NewModelMapperService().InfraReleaseToApplicationModel(release)

		// If previous release exists, set from to there
		var releaseTo infra.Release
		var timeTo *time.Time
		if i != 0 {
			releaseTo = release
			timeTo = &releaseTo.Time
		} else {
			timeTo = nil
		}

		// If next release doesn't exist (UNRELEASED has not next), set to there
		var releaseFrom infra.Release
		var timeFrom *time.Time
		if i != len(infraTruncatedReleases)-1 {
			releaseFrom = infraTruncatedReleases[i+1]
			timeFrom = &releaseFrom.Time
		} else {
			if lastRelease != nil {
				layout := "2006-01-02T15:04:05"
				str := lastRelease.Date + "T00:00:00"
				t, _ := time.Parse(layout, str)
				timeFrom = &t
			} else {
				timeFrom = nil
			}
		}

		// Obtain the tasks between the last release to this one (or to now)
		tasks, err := cgc.releaseProvider.GetTasks(timeFrom, timeTo, defaultBranch)
		if err != nil {
			return nil, err
		}

		fmt.Println("Retrieved " + strconv.Itoa(len(*tasks)) + " tasks for Release " + release.Name + "...")

		// Pass tasks through Task Pipes
		*tasks = cgc.throughTaskPipes(*tasks)

		// Map each task to an application layer model to add it to the release
		for _, task := range *tasks {
			fmt.Println("\t -> " + task.Name + " " + task.Title)
			appTruncatedRelease.Sections[task.Category] =
				append(appTruncatedRelease.Sections[task.Category],
					*services.NewModelMapperService().InfraTaskToApplicationModel(task))
		}

		appTruncatedReleases = append(appTruncatedReleases, appTruncatedRelease)

	}

	var from *time.Time
	var releaseUrl string
	if len(infraTruncatedReleases) == 0 && len(*releases) != 0 {
		from = &(*releases)[len(*releases)-1].Time
		urlPointer, err := cgc.releaseProvider.ReleaseURL(&(*releases)[0].Name, defaultBranch)
		releaseUrl = *urlPointer
		if err != nil {
			return nil, err
		}
	} else if len(infraTruncatedReleases) != 0 {
		from = &infraTruncatedReleases[len(infraTruncatedReleases)-1].Time
		urlPointer, err := cgc.releaseProvider.ReleaseURL(&infraTruncatedReleases[0].Name, defaultBranch)
		releaseUrl = *urlPointer
		if err != nil {
			return nil, err
		}
	} else {
		if lastRelease != nil {
			layout := "2006-01-02T15:04:05"
			str := lastRelease.Date + "T00:00:00"
			t, _ := time.Parse(layout, str)
			from = &t
		} else {
			from = nil
		}
		urlPointer, err := cgc.releaseProvider.ReleaseURL(nil, defaultBranch)
		releaseUrl = *urlPointer
		if err != nil {
			return nil, err
		}
	}

	unreleasedTasks, _ := cgc.taskProvider.GetTasks(from, nil, defaultBranch)

	// Pass tasks through Task Pipes
	*unreleasedTasks = cgc.throughTaskPipes(*unreleasedTasks)

	unreleasedRelease := models.NewRelease("Unreleased", "",
		releaseUrl, false, nil)
	for _, task := range *unreleasedTasks {
		// If the tasks in unreleased were not in previous release, we append it to the final unreleased section
		cm := services2.ChangelogMixer{}
		lastReleaseContainsTask := false
		if lastRelease != nil {
			_, _, lastReleaseContainsTask = cm.ReleaseContainsTask(*lastRelease,
				*services.NewModelMapperService().InfraTaskToApplicationModel(task))
		}

		if !lastReleaseContainsTask {
			category := cgc.releaseProvider.DefineCategory(task)
			unreleasedRelease.Sections[category] = append(unreleasedRelease.Sections[category],
				*services.NewModelMapperService().InfraTaskToApplicationModel(task))
		}
	}

	appTruncatedReleases = append(appTruncatedReleases, *unreleasedRelease)

	helpers.ReverseAny(appTruncatedReleases)
	return &appTruncatedReleases, nil
}

func (cgc *ChangelogGuardianController) GetTask(taskId string) (*models.Task, error) {
	task, err := cgc.taskProvider.GetTask(taskId)
	if err != nil {
		return nil, err
	}

	taskFromProvider := cgc.throughTaskPipes([]infra.Task{*task})

	if len(taskFromProvider) != 0 {
		return services.NewModelMapperService().InfraTaskToApplicationModel(taskFromProvider[0]), nil
	} else {
		return &models.Task{}, nil
	}
}

func (cgc *ChangelogGuardianController) throughReleasePipes(releases []infra.Release) []infra.Release {
	// Reverse to start from the first item
	helpers.ReverseAny(cgc.releasePipes)

	var finalReleases = releases
	for _, releasePipe := range cgc.releasePipes {

		var provisionalReleases []infra.Release
		for _, release := range finalReleases {
			pipeedRelease, _, err := releasePipe.Pipe(&release)
			if err == nil && pipeedRelease != nil {
				provisionalReleases = append(provisionalReleases, *pipeedRelease)
			}
		}
		finalReleases = provisionalReleases
		provisionalReleases = []infra.Release{}
	}
	return finalReleases
}

func (cgc *ChangelogGuardianController) throughTaskPipes(tasks []infra.Task) []infra.Task {
	// Reverse to start from the first item
	helpers.ReverseAny(cgc.taskPipes)

	var finalTasks = tasks
	for _, taskPipe := range cgc.taskPipes {

		var provisionalTasks []infra.Task
		for _, task := range finalTasks {
			pipeedTask, _, err := taskPipe.Pipe(&task)
			if err == nil && pipeedTask != nil {
				provisionalTasks = append(provisionalTasks, *pipeedTask)
			}
		}
		finalTasks = provisionalTasks
		provisionalTasks = []infra.Task{}
	}
	return finalTasks
}
