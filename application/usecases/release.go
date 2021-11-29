package usecases

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/selectors"
	"gitlab.com/aoterocom/changelog-guardian/application/services"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/helpers"
	"strconv"
	"time"
)

func ReleaseCmd(cmd *cobra.Command, args []string) {
	// Update current CHANGELOG to prepare for release, using the regular command
	changelog := RegularCmd(cmd, args)

	argTemplate := cmd.Flag("template").Value.String()
	if argTemplate != "" {
		Settings.Template = argTemplate
	}
	changelogService, err := selectors.ChangelogTemplateSelector(Settings.Template)
	if err != nil {
		panic(err)
	}

	// Load args:
	argPatch, _ := strconv.ParseBool(cmd.Flag("patch").Value.String())
	argMinor, _ := strconv.ParseBool(cmd.Flag("minor").Value.String())
	argMajor, _ := strconv.ParseBool(cmd.Flag("major").Value.String())
	argForce, _ := strconv.ParseBool(cmd.Flag("force").Value.String())
	argVersion := cmd.Flag("version").Value.String()
	argPrerelease := cmd.Flag("pre").Value.String()
	argBuild := cmd.Flag("build").Value.String()

	var versionToBump string
	var nextVersion string
	var lastRelease *models.Release
	semVer := services.NewSemVerService()

	if argForce && argVersion != "" {
		nextVersion = argVersion
		Log.Warnf("\"WARNING: forcing a specific version is not recommended\n")
	} else if len(changelog.Releases) > 1 {
		unreleased := &changelog.Releases[0]
		lastRelease = &changelog.Releases[1]
		versionToBump = lastRelease.Version

		// Calculate the coming up version
		var categories []models.Category
		for key, _ := range unreleased.Sections {
			categories = append(categories, key)
		}

		nextVersion = semVer.CalculateNextVersion(categories, versionToBump)

		if !argForce {
			if argPatch {
				checkVersion := semVer.BumpPatch(versionToBump)
				if nextVersion != checkVersion {
					Log.Errorf("Cannot release a patch version: current tasks imply a minor/major version bump\n")
					return
				}
			}
			if argMinor {
				checkVersion := semVer.BumpMinor(versionToBump)
				if nextVersion != checkVersion {
					Log.Errorf("Cannot release a minor version: no current tasks implying a minor version " +
						"or at least implying breaking changes\n")
					return
				}
			}
			if argMajor {
				checkVersion := semVer.BumpMinor(versionToBump)
				if nextVersion != checkVersion {
					Log.Errorf("Cannot release a major version: no breaking changes found\n")
					return
				}
			}
			if argVersion != "" {
				if nextVersion != argVersion {
					Log.Errorf("Version breaks semver: current version expected would be %s", nextVersion)
					return
				}
			}
		} else {
			if argPatch {
				nextVersion = semVer.BumpPatch(versionToBump)
			}
			if argMinor {
				nextVersion = semVer.BumpMinor(versionToBump)
			}
			if argMajor {
				nextVersion = semVer.BumpMajor(versionToBump)
			}
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

	Log.Infof("Preparing Changelog with Release %s\n", nextVersion)

	// Convert Unreleased to nextVersion and add new UNRELEASE section
	unreleased := &changelog.Releases[0]
	unreleased.Version = nextVersion

	releaseProvider, err := selectors.ProviderSelector(Settings.ReleaseProvider)
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
	newUnreleasedURL, err := (*releaseProvider).ReleaseURL(&unreleased.Version, Settings.DevelopBranch)
	changelog.Releases = append(changelog.Releases, *models.NewRelease("Unreleased", "", *newUnreleasedURL, false, nil))
	helpers.ReverseAny(changelog.Releases)

	argOutputTemplate := cmd.Flag("output-template").Value.String()
	if argOutputTemplate != "" {
		changelogService, err = selectors.ChangelogTemplateSelector(argOutputTemplate)
	}
	err = (*changelogService).SaveChangelog(*changelog, Settings.ChangelogPath)
	if err != nil {
		panic(err)
	}

	Log.Infof("Changelog saved\n")

	if Log.GetLevel() <= log.ErrorLevel {
		fmt.Printf("%s", nextVersion)
	}
}
