package usecases

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/selectors"
	"gitlab.com/aoterocom/changelog-guardian/application/services"
	. "gitlab.com/aoterocom/changelog-guardian/config"
)

func YankCmd(cmd *cobra.Command, args []string) {

	argTemplate := cmd.Flag("template").Value.String()
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

	var versionToYank string
	if argVersion == "" {
		if len(localChangelog.Releases) > 1 {
			lastRelease := &localChangelog.Releases[1]
			if lastRelease.Yanked {
				fmt.Println("Skipping: version " + lastRelease.Version + " is already yanked.")
				return
			}
			versionToYank = lastRelease.Version
		}
	} else {
		versionToYank = argVersion
	}

	for i, release := range localChangelog.Releases {
		if release.Version == versionToYank {
			fmt.Println("Yanking version " + release.Version)
			cm := services.NewChangelogMixer()

			fmt.Println("Look up to find later release...")
			j := i - 1
			for localChangelog.Releases[j].Yanked || i < 0 {
				j--
			}

			fmt.Println("Moving tasks to " + localChangelog.Releases[j].Version)
			localChangelog.Releases[j].Sections = cm.MergeReleases(localChangelog.Releases[j], localChangelog.Releases[i]).Sections
			localChangelog.Releases[i].Sections = make(map[models.Category][]models.Task)
			localChangelog.Releases[i].Yanked = true
			break
		}
	}

	err = (*changelogService).SaveChangelog(*localChangelog, Settings.ChangelogPath)
	if err != nil {
		panic(err)
	}
}
