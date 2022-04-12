package pipes

import (
	settings "gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/controller/services"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"regexp"
)

type JiraTasksPipe struct {
}

func NewJiraTasksPipe() *InclusionsExclusionsTasksPipe {
	return &InclusionsExclusionsTasksPipe{}
}

func (tf *JiraTasksPipe) Filter(task *infra.Task) (*infra.Task, bool, error) {
	var (
		regex       = settings.Settings.TasksPipesCfg.Jira.REGEX
		jiraService = services.NewJiraService()
	)

	rg := regexp.MustCompile(regex)
	match := rg.FindStringSubmatch(task.Title)
	paramsMap := make(map[string]string)
	for i, name := range rg.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	if len(match) != 0 {
		grabbedTask, err := jiraService.GetTask(paramsMap["key"])
		if err != nil {
			return nil, true, err
		}
		return &grabbedTask, true, nil
	}

	return nil, true, nil
}
