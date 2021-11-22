package controllers

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"gitlab.com/aoterocom/changelog-guardian/application/helpers"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/controller/interfaces"
	"gitlab.com/aoterocom/changelog-guardian/controller/services"
	infraInterfaces "gitlab.com/aoterocom/changelog-guardian/infrastructure/interfaces"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"os"
	"strconv"
	"time"
)

type ChangelogGuardianController struct {
	releaseProvider infraInterfaces.Provider
	taskProvider    infraInterfaces.Provider
	releaseFilters  []interfaces.ReleaseFilter
	taskFilters     []interfaces.TaskFilter
}

func NewChangelogGuardianController(releaseProvider infraInterfaces.Provider, taskProvider infraInterfaces.Provider, releaseFiltersStr []string, taskFiltersStr []string) (*ChangelogGuardianController, error) {
	// reversing task and release filters to start iterating over the last
	helpers.ReverseAny(releaseFiltersStr)
	helpers.ReverseAny(taskFiltersStr)

	var releaseFilters []interfaces.ReleaseFilter
	for _, filterStr := range releaseFiltersStr {
		filter, err := services.ReleaseFilterSelector(filterStr)
		if err != nil {
			return nil, err
		}
		releaseFilters = append(releaseFilters, *filter)
	}

	var taskFilters []interfaces.TaskFilter
	for _, filterStr := range taskFiltersStr {
		filter, err := services.TaskFilterSelector(filterStr)
		if err != nil {
			return nil, err
		}
		taskFilters = append(taskFilters, *filter)
	}

	return &ChangelogGuardianController{releaseProvider: releaseProvider, taskProvider: taskProvider, releaseFilters: releaseFilters, taskFilters: taskFilters}, nil
}

func (cgc *ChangelogGuardianController) CetFilledReleasesFromInfra(lastRelease *models.Release, mainBranch string, defaultBranch string) (*[]models.Release, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	r, err := git.PlainOpen(cwd)
	if err != nil {
		return nil, err
	}
	remotes, _ := r.Remotes()
	currentGitBAseUrl := remotes[0].Config().URLs[0]

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

	// Pass releases through Release Filters
	infraTruncatedReleases = cgc.throughReleaseFilters(infraTruncatedReleases)

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
			mainBranch = "develop"
			timeFrom = nil
		}

		// Obtain the tasks between the last release to this one (or to now)
		tasks, err := cgc.releaseProvider.GetTasks(timeFrom, timeTo, mainBranch)
		if err != nil {
			return nil, err
		}

		fmt.Println("Retrieved " + strconv.Itoa(len(*tasks)) + " tasks for Release " + release.Name + "...")

		// Pass tasks through Task Filters
		*tasks = cgc.throughTaskFilters(*tasks)

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
		releaseUrl = cgc.releaseProvider.ReleaseURL(currentGitBAseUrl, &(*releases)[0].Name, defaultBranch)
	} else if len(infraTruncatedReleases) != 0 {
		from = &infraTruncatedReleases[len(infraTruncatedReleases)-1].Time
		releaseUrl = cgc.releaseProvider.ReleaseURL(currentGitBAseUrl, &infraTruncatedReleases[0].Name, defaultBranch)
	} else {
		from = nil
		releaseUrl = cgc.releaseProvider.ReleaseURL(currentGitBAseUrl, nil, defaultBranch)
	}

	unreleasedTasks, _ := cgc.taskProvider.GetTasks(from, nil, defaultBranch)

	// Pass tasks through Task Filters
	*unreleasedTasks = cgc.throughTaskFilters(*unreleasedTasks)

	unreleasedRelease := models.NewRelease("UNRELEASED", "",
		releaseUrl, false, nil)
	for _, task := range *unreleasedTasks {
		category := cgc.releaseProvider.DefineCategory(task)
		unreleasedRelease.Sections[category] = append(unreleasedRelease.Sections[category],
			*services.NewModelMapperService().InfraTaskToApplicationModel(task))
	}

	appTruncatedReleases = append(appTruncatedReleases, *unreleasedRelease)

	helpers.ReverseAny(appTruncatedReleases)
	return &appTruncatedReleases, nil
}

func (cgc *ChangelogGuardianController) throughReleaseFilters(releases []infra.Release) []infra.Release {
	// Reverse to start from the first item
	helpers.ReverseAny(cgc.releaseFilters)

	var finalReleases = releases
	for _, releaseFilter := range cgc.releaseFilters {

		var provisionalReleases []infra.Release
		for _, release := range finalReleases {
			filteredRelease, _, err := releaseFilter.Filter(&release)
			if err == nil && filteredRelease != nil {
				provisionalReleases = append(provisionalReleases, *filteredRelease)
			}
		}
		finalReleases = provisionalReleases
		provisionalReleases = []infra.Release{}
	}
	return finalReleases
}

func (cgc *ChangelogGuardianController) throughTaskFilters(tasks []infra.Task) []infra.Task {
	// Reverse to start from the first item
	helpers.ReverseAny(cgc.taskFilters)

	var finalTasks = tasks
	for _, taskFilter := range cgc.taskFilters {

		var provisionalTasks []infra.Task
		for _, task := range finalTasks {
			filteredTask, _, err := taskFilter.Filter(&task)
			if err == nil && filteredTask != nil {
				provisionalTasks = append(provisionalTasks, *filteredTask)
			}
		}
		finalTasks = provisionalTasks
		provisionalTasks = []infra.Task{}
	}
	return finalTasks
}
