package usecases

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/services"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/controller/controllers"
)

func RegularCmd() *models.Changelog {

	releaseProvider, err := services.ProviderSelector(Settings.ReleaseProvider)
	if err != nil {
		panic(err)
	}

	tasksProvider, err := services.ProviderSelector(Settings.TasksProvider)
	if err != nil {
		panic(err)
	}

	cgController, err := controllers.NewChangelogGuardianController(*releaseProvider, *tasksProvider, Settings.ReleaseFilters, Settings.TaskFilters)
	if err != nil {
		panic(err)
	}

	localChangelog, err := services.ParseChangelog(Settings.ChangelogPath)
	if err != nil && err == errors.Errorf("open : no such file or directory") {
		panic(err)
	}

	var lastRelease *models.Release

	if localChangelog != nil {
		if len(localChangelog.Releases) > 1 {
			lastRelease = &localChangelog.Releases[1]
		}
	} else {
		lastRelease = nil
	}

	releases, err := cgController.GetFilledReleasesFromInfra(lastRelease, Settings.MainBranch, Settings.DevelopBranch)
	if err != nil {
		panic(err)
	}
	retrievedChangelog := models.NewChangelog(*releases)

	if localChangelog != nil {
		changelogMixer := services.NewChangelogMixer()

		// Before merge the changelogs, we need to ensure that the retrieved unreleased section doesn't contain a task
		// already used in the local changelog
		if localChangelog != nil && len(localChangelog.Releases) > 0 {
			retrievedUnreleased := retrievedChangelog.Releases[0]
			localChangelogExceptUnreleased := localChangelog.Releases[:len(localChangelog.Releases)-1]
			for sec, sectionTasks := range retrievedUnreleased.Sections {
				for i, task := range sectionTasks {
					_, _, exists := changelogMixer.ChangelogContainsTask(*models.NewChangelog(localChangelogExceptUnreleased), task)
					if exists {
						retrievedUnreleased.Sections[sec] = remove(sectionTasks, i)
					}
				}
			}
			retrievedChangelog.Releases[0] = retrievedUnreleased
		}

		mergedChangelog := changelogMixer.MergeChangelogs(*localChangelog, *retrievedChangelog)
		retrievedChangelog = &mergedChangelog
	}

	err = retrievedChangelog.Save(Settings.ChangelogPath)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done")

	return retrievedChangelog
}

func remove(slice []models.Task, s int) []models.Task {
	return append(slice[:s], slice[s+1:]...)
}
