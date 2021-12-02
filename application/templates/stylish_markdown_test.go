package templates

import (
	_ "embed"
	"github.com/magiconair/properties/assert"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	settings "gitlab.com/aoterocom/changelog-guardian/config"
	"io/ioutil"
	"os"
	"testing"
)

const OkSTask = "- ✨ [TASK1](https://gitlab.com/aoterocom/changelog-guardian/issues/task1) This is the task 1 ([@aoterocom](https://gitlab.com/aoterocom/))"

func TestSTask(t *testing.T) {
	task := models.NewEmptyTask()

	changelogService := StylishMarkDownChangelogService{}
	parsedTaskPtr := changelogService.parseTask(OkSTask, "")
	task = &parsedTaskPtr
	assert.Equal(t, changelogService.TaskToString(*task, models.ADDED), OkSTask)
}

func TestSChangelogParsing(t *testing.T) {
	settings.Settings.RetrieveSettingsFromFile("../../config/settings.yml")

	cwd, _ := os.Getwd()
	changelogService := StylishMarkDownChangelogService{}
	pathToChangelog := cwd + "/" + "resources/stylish_markdown_CHANGELOG.md"

	changelog, err := changelogService.Parse(pathToChangelog)
	if err != nil {
		panic(err)
	}

	fullChangelog, err := ioutil.ReadFile(pathToChangelog)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, changelogService.String(*changelog), string(fullChangelog))
}
