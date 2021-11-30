package usecases

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/selectors"
	"gitlab.com/aoterocom/changelog-guardian/application/services"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	"strconv"
)

func CalculateReleaseCmd(cmd *cobra.Command, args []string) {
	// Update current CHANGELOG to prepare for release, using the regular command
	Log.SetLevel(log.ErrorLevel)

	Log.Debugf("Preparing execution...\n")
	argTemplate := cmd.Flag("template").Value.String()
	if argTemplate != "" {
		Settings.Template = argTemplate
	}

	Log.Debugf("Using %s template\n", Settings.Template)
	changelogService, err := selectors.ChangelogTemplateSelector(Settings.Template)
	if err != nil {
		Log.Fatalf("Error selecting template\n")
	}

	Log.Infof("Retrieving changelog from %s...\n", Settings.ChangelogPath)
	changelog, err := (*changelogService).Parse(Settings.ChangelogPath)
	if err != nil && err == errors.Errorf("open : no such file or directory") {
		Log.Fatalf("Error retrieving changelog file\n")
	}

	// Load args:
	argPatch, _ := strconv.ParseBool(cmd.Flag("patch").Value.String())
	argMinor, _ := strconv.ParseBool(cmd.Flag("minor").Value.String())
	argMajor, _ := strconv.ParseBool(cmd.Flag("major").Value.String())
	argPrerelease := cmd.Flag("pre").Value.String()
	argBuild := cmd.Flag("build").Value.String()

	var versionToBump string
	var nextVersion string
	var lastRelease *models.Release
	semVer := services.NewSemVerService()

	if len(changelog.Releases) > 1 {
		unreleased := &changelog.Releases[0]
		lastRelease = &changelog.Releases[1]
		versionToBump = lastRelease.Version

		// Calculate the coming up version
		var categories []models.Category
		for key, _ := range unreleased.Sections {
			categories = append(categories, key)
		}

		nextVersion = semVer.CalculateNextVersion(categories, versionToBump)

		if argPatch {
			nextVersion = semVer.BumpPatch(versionToBump)
		}
		if argMinor {
			nextVersion = semVer.BumpMinor(versionToBump)
		}
		if argMajor {
			nextVersion = semVer.BumpMajor(versionToBump)
		}

	} else {
		nextVersion = Settings.InitialVersion
	}

	// Add pre-release and build strings if provided
	if argPrerelease != "" {
		nextVersion = nextVersion + "-" + argPrerelease
	}
	if argBuild != "" {
		nextVersion = nextVersion + "+" + argBuild
	}

	if !semVer.IsSemverValid(nextVersion) {
		Log.Errorf("Provider build metadata or pre-release breaks Semver. "+
			"The result version string %s is not valid", nextVersion)
		return
	}

	fmt.Printf("%s", nextVersion)
}
