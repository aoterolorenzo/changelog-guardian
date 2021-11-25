package themes

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"io/ioutil"
	"os"
	"testing"
)

const OkTask = "- [TASK1](https://gitlab.com/aoterocom/changelog-guardian/issues/task1) This is the task 1 ([@aoterocom](https://gitlab.com/aoterocom/))"

func TestTask(t *testing.T) {
	task := models.NewEmptyTask()

	changelogService := MarkDownChangelogService{}
	parsedTaskPtr := changelogService.parseTask(OkTask, "")
	task = &parsedTaskPtr
	fmt.Println(task)
	assert.Equal(t, task.String(), OkTask)
}

func TestChangelogParsing(t *testing.T) {
	cwd, _ := os.Getwd()
	changelogService := MarkDownChangelogService{}
	pathToChangelog := cwd + "/" + "resources/CHANGELOG.md"

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
