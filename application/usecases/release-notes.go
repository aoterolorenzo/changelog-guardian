package usecases

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/aoterocom/changelog-guardian/application/helpers"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/selectors"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	"strconv"
	"strings"
)

func ReleaseNotesCmd(cmd *cobra.Command, args []string) {

	Log.Debugf("Preparing execution...\n")
	argTemplate := cmd.Flag("template").Value.String()
	argOutputFile := cmd.Flag("output-file").Value.String()
	argEcho, _ := strconv.ParseBool(cmd.Flag("echo").Value.String())

	if argEcho {
		Log.SetLevel(logrus.ErrorLevel)
	}

	if argTemplate != "" {
		Settings.Template = argTemplate
	}
	Log.Debugf("Using %s template\n", Settings.Template)
	changelogService, err := selectors.ChangelogTemplateSelector(Settings.Template)
	if err != nil {
		panic(err)
	}

	Log.Infof("Retrieving changelog from %s...\n", Settings.ChangelogPath)
	localChangelog, err := (*changelogService).Parse(Settings.ChangelogPath)
	if err != nil && err == errors.Errorf("open : no such file or directory") {
		Log.Errorf("Changelog not found at %s\n", Settings.ChangelogPath)
		return
	}

	argVersion := cmd.Flag("version").Value.String()

	Log.Infof("Preparing Release Notes...\n")
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

			argOutputTemplate := cmd.Flag("output-template").Value.String()
			if argOutputTemplate != "" {
				changelogService, err = selectors.ChangelogTemplateSelector(argOutputTemplate)
			}

			releaseNotes := (*changelogService).NudeChangelogString(*changelog)
			if releaseNotes != "" {
				// Truncates initial and ending line break it exists
				releaseNotes = strings.TrimPrefix(releaseNotes, "\n")
				releaseNotes = strings.TrimSuffix(releaseNotes, "\n")
			}

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

				Log.Infof("Release Notes saved on %s\n", Settings.ReleaseNotesPath)
			}
			return
		}
	}

	Log.Errorf("Version not found\n")
}
