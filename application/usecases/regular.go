package usecases

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/selectors"
	"gitlab.com/aoterocom/changelog-guardian/application/services"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/controller/controllers"
	infraSelectors "gitlab.com/aoterocom/changelog-guardian/infrastructure/selectors"
	"strconv"
)

func RegularCmd(cmd *cobra.Command, args []string) *models.Changelog {

	Log.Debugf("Preparing execution...\n")
	argTemplate := cmd.Flag("template").Value.String()
	if argTemplate != "" {
		Settings.Template = argTemplate
	}
	Log.Debugf("Using %s template\n", Settings.Template)
	changelogService, err := selectors.ChangelogTemplateSelector(Settings.Template)
	if err != nil {
		Log.WithError(err).Fatalf("Error selecting template\n")
	}

	Log.Debugf("Release provider: %s\n", Settings.ReleaseProvider)
	releaseProvider, err := infraSelectors.ProviderSelector(Settings.ReleaseProvider)
	if err != nil {
		Log.WithError(err).Fatalf("Error selecting release provider\n")
	}

	Log.Debugf("Tasks provider: %s\n", Settings.TasksProvider)
	tasksProvider, err := infraSelectors.ProviderSelector(Settings.TasksProvider)
	if err != nil {
		Log.WithError(err).Fatalf("Error selecting tasks provider\n")
	}

	cgController, err := controllers.NewChangelogGuardianController(*releaseProvider, *tasksProvider, Settings.ReleasePipes, Settings.TasksPipes)
	if err != nil {
		Log.WithError(err).Fatalf("Error creating controller\n")
	}

	Log.Infof("Retrieving changelog from %s...\n", Settings.ChangelogPath)
	localChangelog, err := (*changelogService).Parse(Settings.ChangelogPath)
	if err != nil && err == errors.Errorf("open : no such file or directory") {
		Log.WithError(err).Fatalf("Error retrieving changelog file\n")
	}

	var lastRelease *models.Release

	if localChangelog != nil {
		if len(localChangelog.Releases) > 1 {
			lastRelease = &localChangelog.Releases[1]
			Log.Infof("Retrieving releases and tasks data from Release %s...\n", lastRelease.Version)
		}
	} else {
		lastRelease = nil
		Log.Infof("Retrieving release and tasks data...\n")
	}

	Log.Infof("Preparing Changelog...\n")

	releases, err := cgController.GetFilledReleasesFromInfra(lastRelease, Settings.MainBranch, Settings.DevelopBranch)
	if err != nil && err.Error() == "repository does not exist" {
		Log.WithError(err).Fatalf("No git repository found on this path")
	} else if err != nil {
		Log.WithError(err).Fatalf(err.Error())
	}
	retrievedChangelog := models.NewChangelog()
	retrievedChangelog.Releases = *releases

	if localChangelog != nil {
		changelogMixer := services.NewChangelogMixer()

		// Before merge the changelogs, we need to ensure that the retrieved unreleased section doesn't contain a task
		// already used in the local changelog
		if localChangelog != nil && len(localChangelog.Releases) > 0 {
			retrievedUnreleased := retrievedChangelog.Releases[0]
			localChangelogExceptUnreleased := localChangelog.Releases[:len(localChangelog.Releases)-1]
			for sec, sectionTasks := range localChangelog.Releases[0].Sections {
				// If --no-dup set, we just remove duplicates
				argNoDup, _ := strconv.ParseBool(cmd.Flag("no-dups").Value.String())
				if argNoDup {
					// Create a map to store unique IDs
					uniqueIDs := make(map[string]bool)
					// Create a slice to store non-duplicate tasks
					nonDuplicates := make([]models.Task, 0)
					for _, task := range sectionTasks {
						if _, exists := uniqueIDs[task.ID]; !exists {
							nonDuplicates = append(nonDuplicates, task)
							uniqueIDs[task.ID] = true
						}
					}
					localChangelog.Releases[0].Sections[sec] = nonDuplicates
				}
			}

			for sec, sectionTasks := range retrievedUnreleased.Sections {
				for i, task := range sectionTasks {
					provChangelog := *models.NewChangelog()
					provChangelog.Releases = localChangelogExceptUnreleased
					_, _, exists := changelogMixer.ChangelogContainsTask(provChangelog, task)
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

	argOutputTemplate := cmd.Flag("output-template").Value.String()
	if argOutputTemplate != "" {
		changelogService, err = selectors.ChangelogTemplateSelector(argOutputTemplate)
	}
	err = (*changelogService).SaveChangelog(*retrievedChangelog, Settings.ChangelogPath)
	if err != nil {
		Log.WithError(err).Fatalf("Error saving changelog file\n")
	}

	Log.Infof("Changelog saved\n")

	return retrievedChangelog
}

func remove(slice []models.Task, s int) []models.Task {
	return append(slice[:s], slice[s+1:]...)
}
