package usecases

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/aoterocom/changelog-guardian/application/helpers"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/selectors"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	"strconv"
)

func ReleaseNotesCmd(cmd *cobra.Command, args []string) {

	argTemplate := cmd.Flag("template").Value.String()
	argOutputFile := cmd.Flag("output-file").Value.String()
	argEcho, _ := strconv.ParseBool(cmd.Flag("echo").Value.String())

	if argTemplate != "" {
		Settings.Style = argTemplate
	}
	changelogService, err := selectors.ChangelogTemplateSelector(Settings.Style)
	if err != nil {
		panic(err)
	}

	localChangelog, err := (*changelogService).Parse(Settings.ChangelogPath)
	if err != nil && err == errors.Errorf("open : no such file or directory") {
		panic(err)
	}

	argVersion := cmd.Flag("version").Value.String()

	var targetVersion string
	if argVersion == "" {
		if len(localChangelog.Releases) > 1 {
			lastRelease := &localChangelog.Releases[1]
			targetVersion = lastRelease.Version
		}
	} else {
		targetVersion = argVersion
	}

	for _, release := range localChangelog.Releases {
		if release.Version == targetVersion {
			changelog := models.NewChangelog()
			changelog.Releases = []models.Release{release}

			releaseNotes := (*changelogService).NudeChangelogString(*changelog)
			if argEcho {
				fmt.Println(releaseNotes)
			} else {

				if argOutputFile != "" {
					Settings.ReleaseNotesPath = argOutputFile
				}

				err = helpers.SaveStringToFile(Settings.ReleaseNotesPath, releaseNotes)
				if err != nil {
					panic(err)
				}

				fmt.Println("Release Notes written on " + Settings.ReleaseNotesPath)
			}
			return
		}
	}

	fmt.Println("Version not found.")
}
