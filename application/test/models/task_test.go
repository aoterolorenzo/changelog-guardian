package models

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/services"
	"io/ioutil"
	"testing"
)

const OkTask = "- [TASK1](https://gitlab.com/aoterocom/changelog-guardian/issues/task1) This is the task 1 ([@aoterocom](https://gitlab.com/aoterocom/))"

func TestTask(t *testing.T) {
	task := models.NewEmptyTask()
	task = services.ParseTask(OkTask)
	fmt.Println(task)
	assert.Equal(t, task.String(), OkTask)
}

func TestChangelogParsing(t *testing.T) {

	pathToChangelog := "/Users/alberto/GolandProyects/CLogger/test/models/resources/CHANGELOG.md"
	changelog, err := services.ParseChangelog(pathToChangelog)
	if err != nil {
		panic(err)
	}
	fullChangelog, err := ioutil.ReadFile(pathToChangelog)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, changelog.String(), string(fullChangelog))
}
