package usecases

import (
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/selectors"
	"gitlab.com/aoterocom/changelog-guardian/application/services"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/controller/controllers"
	"gitlab.com/aoterocom/changelog-guardian/helpers"
	selectors2 "gitlab.com/aoterocom/changelog-guardian/infrastructure/selectors"
	"strconv"
)

func InsertCmd(cmd *cobra.Command, args []string) {

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
	// Load args:
	argTitle := cmd.Flag("title").Value.String()
	argId := cmd.Flag("id").Value.String()
	argLink := cmd.Flag("link").Value.String()
	argAuthor := cmd.Flag("author").Value.String()
	argAuthorLink := cmd.Flag("authorLink").Value.String()
	argCategory := cmd.Flag("category").Value.String()
	argSkipAutocompletion, _ := strconv.ParseBool(cmd.Flag("skip-autocompletion").Value.String())

	if argId == "" {
		Log.Errorf("--id argument is mandatory\n")
	}

	var taskFromProvider = &models.Task{}
	if !argSkipAutocompletion {
		Log.Debugf("Release provider: %s\n", Settings.ReleaseProvider)
		tasksProvider, err := selectors2.ProviderSelector(Settings.TasksProvider)
		if err != nil {
			Log.Fatalf("Error selecting release provider\n")
		}

		cgController, err := controllers.NewChangelogGuardianController(nil, *tasksProvider, Settings.ReleasePipes, Settings.TasksPipes)
		if err != nil {
			Log.Fatalf("Error creating controller\n")
		}

		Log.Debugf("Retrieving task info from provider\n")
		taskFromProvider, err = cgController.GetTask(argId)
		if err != nil {
			return
		}

	}

	Log.Infof("Retrieving changelog from %s...\n", Settings.ChangelogPath)
	localChangelog, err := (*changelogService).Parse(Settings.ChangelogPath)
	if err != nil && err == errors.Errorf("open : no such file or directory") {
		Log.Fatalf("Error retrieving changelog file\n")
	}

	task := models.NewTask(argId, argId, argLink, argTitle, argAuthor, argAuthorLink, models.Category(argCategory))
	err = mergo.Merge(task, taskFromProvider)
	if err != nil {
		return
	}

	cm := services.NewChangelogMixer()
	_, _, exists := cm.ReleaseContainsTask(localChangelog.Releases[0], *task)

	if !exists {
		helpers.ReverseAny(localChangelog.Releases[0].Sections[models.Category(argCategory)])
		localChangelog.Releases[0].Sections[models.Category(argCategory)] =
			append(localChangelog.Releases[0].Sections[models.Category(argCategory)], *task)
		Log.Infof("Task successfully added\n")
		helpers.ReverseAny(localChangelog.Releases[0].Sections[models.Category(argCategory)])
	} else {
		Log.Infof("Task  %s is already present on the CHANGELOG. Skipping...\n", task.ID)
		return
	}

	argOutputTemplate := cmd.Flag("output-template").Value.String()
	if argOutputTemplate != "" {
		changelogService, err = selectors.ChangelogTemplateSelector(argOutputTemplate)
	}
	err = (*changelogService).SaveChangelog(*localChangelog, Settings.ChangelogPath)
	if err != nil {
		return
	}
	Log.Infof("Changelog saved\n")
}
