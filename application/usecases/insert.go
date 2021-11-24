package usecases

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/aoterocom/changelog-guardian/application/helpers"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/services"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/controller/controllers"
	"strconv"
)

func InsertCmd(cmd *cobra.Command, args []string) {

	// Load args:
	argTitle := cmd.Flag("title").Value.String()
	argId := cmd.Flag("id").Value.String()
	argLink := cmd.Flag("link").Value.String()
	argAuthor := cmd.Flag("author").Value.String()
	argAuthorLink := cmd.Flag("authorLink").Value.String()
	argCategory := cmd.Flag("category").Value.String()
	argSkipAutocompletion, _ := strconv.ParseBool(cmd.Flag("skip-autocompletion").Value.String())

	if argId == "" {
		fmt.Println("--id argument is mandatory")
	}

	var taskFromProvider = &models.Task{}
	if !argSkipAutocompletion {
		tasksProvider, err := services.ProviderSelector(Settings.TasksProvider)
		if err != nil {
			panic(err)
		}

		cgController, err := controllers.NewChangelogGuardianController(nil, *tasksProvider, Settings.ReleaseFilters, Settings.TaskFilters)
		if err != nil {
			panic(err)
		}

		taskFromProvider, err = cgController.GetTask(argId)
		if err != nil {
			return
		}

	}

	localChangelog, err := services.ParseChangelog(Settings.ChangelogPath)
	if err != nil && err == errors.Errorf("open : no such file or directory") {
		panic(err)
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
		helpers.ReverseAny(localChangelog.Releases[0].Sections[models.Category(argCategory)])
	} else {
		fmt.Println("Skipped: Task " + task.ID + " is already present on the CHANGELOG.")
		return
	}

	err = localChangelog.Save(Settings.ChangelogPath)
	if err != nil {
		fmt.Println(err)
		return
	}
}
