package usecases

import (
	"gitlab.com/aoterocom/changelog-guardian/application/helpers"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/services"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	"time"
)

func ReleaseCmd() {
	// Update current CHANGELOG to prepare for release, using the regular command
	changelog := RegularCmd()

	var versionToBump string
	var nextVersion string
	var lastRelease *models.Release
	if len(changelog.Releases) > 1 {
		semVer := services.SemVerService{}
		unreleased := &changelog.Releases[0]
		lastRelease = &changelog.Releases[1]
		versionToBump = lastRelease.Version

		// Calculate the coming up version
		var categories []models.Category
		for key, _ := range unreleased.Sections {
			categories = append(categories, key)
		}

		nextVersion = semVer.CalculateNextVersion(categories, versionToBump)
	} else {
		nextVersion = Settings.InitialVersion
	}

	// Convert UNRELEASED to nextVersion and add new UNRELEASE section
	unreleased := &changelog.Releases[0]
	unreleased.Version = nextVersion

	releaseProvider, err := services.ProviderSelector(Settings.ReleaseProvider)
	if err != nil {
		panic(err)
	}
	var from *string
	if lastRelease == nil {
		from = nil
	} else {
		from = &lastRelease.Version
	}
	url, err := (*releaseProvider).ReleaseURL(from, nextVersion)
	if err != nil {
		panic(err)
	}
	unreleased.Link = *url
	unreleased.Date = time.Now().Format("2006-01-02")

	changelog.Releases[0] = *unreleased

	helpers.ReverseAny(changelog.Releases)
	changelog.Releases = append(changelog.Releases, *models.NewRelease("UNRELEASED", "", "", false, nil))
	helpers.ReverseAny(changelog.Releases)

	err = changelog.Save(Settings.ChangelogPath)
	if err != nil {
		panic(err)
	}
}
