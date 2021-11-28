package usecases

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/selectors"
	"gitlab.com/aoterocom/changelog-guardian/application/services"
	. "gitlab.com/aoterocom/changelog-guardian/config"
)

func YankCmd(cmd *cobra.Command, args []string) {

	Log.Debugf("Preparing execution...\n")
	argTemplate := cmd.Flag("template").Value.String()
	if argTemplate != "" {
		Settings.Template = argTemplate
	}
	Log.Debugf("Using %s template\n", Settings.Template)
	changelogService, err := selectors.ChangelogTemplateSelector(Settings.Template)
	if err != nil {
		panic(err)
	}

	localChangelog, err := (*changelogService).Parse(Settings.ChangelogPath)
	if err != nil && err == errors.Errorf("open : no such file or directory") {
		Log.Errorf("Changelog not found at %s\n", Settings.ChangelogPath)
		return
	}

	argVersion := cmd.Flag("version").Value.String()

	var versionToYank string
	if argVersion == "" {
		if len(localChangelog.Releases) > 1 {
			lastRelease := &localChangelog.Releases[1]
			if lastRelease.Yanked {
				Log.Errorf("Skipping: Release %s has been already yanked\n", lastRelease.Version)
				return
			}
			versionToYank = lastRelease.Version
		}
	} else {
		versionToYank = argVersion
	}

	for i, release := range localChangelog.Releases {
		if release.Version == versionToYank {

			Log.Infof("Yanking Release %s...\n", release.Version)
			cm := services.NewChangelogMixer()

			Log.Debugf("Searching for next releases...\n")
			j := i - 1
			for localChangelog.Releases[j].Yanked || i < 0 {
				j--
			}

			Log.Debugf("Moving tasks to Release %s...\n", localChangelog.Releases[j].Version)
			localChangelog.Releases[j].Sections = cm.MergeReleases(localChangelog.Releases[j], localChangelog.Releases[i]).Sections
			localChangelog.Releases[i].Sections = make(map[models.Category][]models.Task)
			localChangelog.Releases[i].Yanked = true
			Log.Infof("Release %s succesfully yanked. Tasks moved to closest Release/Section (%s)...\n", localChangelog.Releases[i].Version, localChangelog.Releases[j].Version)
			break
		}
	}

	argOutputTemplate := cmd.Flag("output-template").Value.String()
	if argOutputTemplate != "" {
		changelogService, err = selectors.ChangelogTemplateSelector(argOutputTemplate)
	}
	err = (*changelogService).SaveChangelog(*localChangelog, Settings.ChangelogPath)
	if err != nil {
		panic(err)
	}

	Log.Infof("Changelog saved\n")
}
