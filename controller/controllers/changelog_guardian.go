package controllers

import (
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	services2 "gitlab.com/aoterocom/changelog-guardian/application/services"
	"gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/controller/interfaces"
	"gitlab.com/aoterocom/changelog-guardian/controller/selectors"
	"gitlab.com/aoterocom/changelog-guardian/controller/services"
	"gitlab.com/aoterocom/changelog-guardian/helpers"
	infraInterfaces "gitlab.com/aoterocom/changelog-guardian/infrastructure/interfaces"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"strconv"
	"time"
)

type ChangelogGuardianController struct {
	releaseProvider infraInterfaces.Provider
	taskProvider    infraInterfaces.Provider
	releasePipes    []interfaces.ReleasePipe
	tasksPipes      []interfaces.TasksPipe
}

func NewChangelogGuardianController(releaseProvider infraInterfaces.Provider, taskProvider infraInterfaces.Provider, releasePipesStr []string, tasksPipesStr []string) (*ChangelogGuardianController, error) {

	var releasePipes []interfaces.ReleasePipe
	for _, pipeStr := range releasePipesStr {
		pipe, err := selectors.ReleasePipeSelector(pipeStr)
		if err != nil {
			return nil, err
		}
		releasePipes = append(releasePipes, *pipe)
	}

	var tasksPipes []interfaces.TasksPipe
	for _, pipeStr := range tasksPipesStr {
		pipe, err := selectors.TasksPipeSelector(pipeStr)
		if err != nil {
			return nil, err
		}
		tasksPipes = append(tasksPipes, *pipe)
	}

	return &ChangelogGuardianController{releaseProvider: releaseProvider, taskProvider: taskProvider, releasePipes: releasePipes, tasksPipes: tasksPipes}, nil
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

	settings.Log.Debugf("Found %d releases\n", len(infraTruncatedReleases))

	// Pass releases through Release Pipes
	infraTruncatedReleases = cgc.throughReleasePipes(infraTruncatedReleases)

	// For each release obtained from infra layer
	for i, release := range infraTruncatedReleases {

		// Map the release to application layer model
		appTruncatedRelease := *services.NewModelMapperService().InfraReleaseToApplicationModel(release)

		// If release is not the first one, 'from' the previous
		var releaseFrom infra.Release
		var timeFrom *time.Time
		if i != 0 {
			releaseFrom = infraTruncatedReleases[i-1]
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

		// Always search until ('to') the release moment
		var releaseTo infra.Release
		var timeTo *time.Time
		releaseTo = release
		timeTo = &releaseTo.Time

		// Obtain the tasks between the last release to this one (or to now)
		tasks, err := cgc.releaseProvider.GetTasks(timeFrom, timeTo, defaultBranch)
		if err != nil {
			return nil, err
		}

		settings.Log.Debugf("Retrieved %s tasks for Release %s\n", strconv.Itoa(len(*tasks)), release.Name)
		// Pass tasks through Task Pipes
		*tasks = cgc.throughTasksPipes(*tasks)

		// Map each task to an application layer model to add it to the release
		for _, task := range *tasks {
			settings.Log.Debugf("-> %s %s\n", task.Name, task.Title)
			appTruncatedRelease.Sections[task.Category] =
				append(appTruncatedRelease.Sections[task.Category],
					*services.NewModelMapperService().InfraTaskToApplicationModel(task))
		}

		appTruncatedReleases = append(appTruncatedReleases, appTruncatedRelease)

	}

	var from *time.Time
	var releaseUrl string
	if len(infraTruncatedReleases) == 0 && len(*releases) != 0 {
		from = &(*releases)[0].Time
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
	settings.Log.Debugf("Retieved %d new unreleased tasks\n", len(*unreleasedTasks))

	// Pass tasks through Task Pipes
	*unreleasedTasks = cgc.throughTasksPipes(*unreleasedTasks)

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
			settings.Log.Debugf("-> %s %s\n", task.Name, task.Title)
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

	taskFromProvider := cgc.throughTasksPipes([]infra.Task{*task})

	if len(taskFromProvider) != 0 {
		return services.NewModelMapperService().InfraTaskToApplicationModel(taskFromProvider[0]), nil
	} else {
		return &models.Task{}, nil
	}
}

func (cgc *ChangelogGuardianController) throughReleasePipes(releases []infra.Release) []infra.Release {

	var finalReleases = releases
	for _, releasePipe := range cgc.releasePipes {

		var provisionalReleases []infra.Release
		for _, release := range finalReleases {
			pipeedRelease, _, err := releasePipe.Filter(&release)
			if err == nil && pipeedRelease != nil {
				provisionalReleases = append(provisionalReleases, *pipeedRelease)
			}
		}
		finalReleases = provisionalReleases
		provisionalReleases = []infra.Release{}
	}
	return finalReleases
}

func (cgc *ChangelogGuardianController) throughTasksPipes(tasks []infra.Task) []infra.Task {

	var finalTasks = tasks
	for _, tasksPipe := range cgc.tasksPipes {

		var provisionalTasks []infra.Task
		for _, task := range finalTasks {
			pipeedTask, _, err := tasksPipe.Filter(&task)
			if err == nil && pipeedTask != nil {
				provisionalTasks = append(provisionalTasks, *pipeedTask)
			}
		}
		finalTasks = provisionalTasks
		provisionalTasks = []infra.Task{}
	}
	return finalTasks
}
