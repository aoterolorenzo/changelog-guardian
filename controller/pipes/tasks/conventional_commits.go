package pipes

import (
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	settings "gitlab.com/aoterocom/changelog-guardian/config"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"regexp"
)

type ConventionalCommitsTasksPipe struct {
}

func NewConventionalCommitsTasksPipe() *ConventionalCommitsTasksPipe {
	return &ConventionalCommitsTasksPipe{}
}

func (tf *ConventionalCommitsTasksPipe) Filter(task *infra.Task) (*infra.Task, bool, error) {

	var regex = "^(?P<type>[\\w]+)(?P<scope>\\([\\w\\-\\.]+\\))?(?P<breaking_change_marker>!)?:\\s(?P<description>[^\\n]+)([\\s\\S]*)"
	rg := regexp.MustCompile(regex)
	match := rg.FindStringSubmatch(task.Title)
	paramsMap := make(map[string]string)
	for i, name := range rg.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}

	if len(match) != 0 {
		if paramsMap["breaking_change_marker"] == "!" {
			task.Category = models.BREAKING_CHANGE
		} else {
			for key, val := range settings.Settings.TasksPipesCfg.ConventionalCommits.Categories {
				if paramsMap["type"] == val {
					task.Category = key
					break
				}
			}
		}

		return task, false, nil
	}

	return nil, true, nil
}
