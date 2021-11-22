package usecases

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/services"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/controller/controllers"
	controllerInterfaces "gitlab.com/aoterocom/changelog-guardian/controller/interfaces"
	"gitlab.com/aoterocom/changelog-guardian/infrastructure/interfaces"
	"gitlab.com/aoterocom/changelog-guardian/infrastructure/providers"
)

func Regular() {

	releaseProvider := interfaces.Provider(providers.NewGitlabProvider())
	tasksProvider := interfaces.Provider(providers.NewGitlabProvider())

	cgController := controllers.NewChangelogGuardianController(releaseProvider, tasksProvider, []controllerInterfaces.Middleware{})

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
