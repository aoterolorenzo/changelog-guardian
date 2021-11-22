package usecases

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/services"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/controller/controllers"
)

func Regular() {

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

	releases, err := cgController.CetFilledReleasesFromInfra(lastRelease, Settings.MainBranch, Settings.DevelopBranch)
	if err != nil {
		panic(err)
	}
	retrievedChangelog := models.NewChangelog(*releases)

	if localChangelog != nil {
		changelogMixer := services.NewChangelogMixer()
		mergedChangelog := changelogMixer.MergeChangelogs(*localChangelog, *retrievedChangelog)
		retrievedChangelog = &mergedChangelog
	}

	err = retrievedChangelog.Save(Settings.ChangelogPath)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done")

}
