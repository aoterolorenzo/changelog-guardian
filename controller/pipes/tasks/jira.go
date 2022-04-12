package pipes

import (
	settings "gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/controller/controllers/providers"
	"gitlab.com/aoterocom/changelog-guardian/controller/interfaces"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"regexp"
)

type JiraTasksPipe struct {
	providerController interfaces.ProviderController
}

func NewJiraTasksPipe() *JiraTasksPipe {
	jiraController := providers.NewJiraController()
	return &JiraTasksPipe{
		providerController: interfaces.ProviderController(jiraController),
	}
}

func (tf *JiraTasksPipe) Filter(task *infra.Task) (*infra.Task, bool, error) {

	var regex = settings.Settings.TasksPipesCfg.Jira.REGEX
	rg := regexp.MustCompile(regex)
	match := rg.FindStringSubmatch(task.Title)
	paramsMap := make(map[string]string)
	for i, name := range rg.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	if len(match) != 0 {
		grabbedTask, err := tf.providerController.GetTask(paramsMap["key"])
		if err != nil {
			return nil, true, err
		}
		return grabbedTask, true, nil
	}

	return nil, true, nil
}
