package selectors

import (
	"github.com/pkg/errors"
	"gitlab.com/aoterocom/changelog-guardian/controller/interfaces"
	tasksPipes "gitlab.com/aoterocom/changelog-guardian/controller/pipes/tasks"
)

func TasksPipeSelector(providerStr string) (*interfaces.TasksPipe, error) {
	switch providerStr {
	case "gitlab_resolver":
		prov := interfaces.TasksPipe(tasksPipes.NewGitlabResolverTasksPipe())
		return &prov, nil
	case "natural_language":
		prov := interfaces.TasksPipe(tasksPipes.NewNaturalLanguageTasksPipe())
		return &prov, nil
	case "conventional_commits":
		prov := interfaces.TasksPipe(tasksPipes.NewConventionalCommitsTasksPipe())
		return &prov, nil
	case "inclusions_exclusions":
		prov := interfaces.TasksPipe(tasksPipes.NewInclusionsExclusionsTasksPipe())
		return &prov, nil
	case "jira":
		prov := interfaces.TasksPipe(tasksPipes.NewJiraTasksPipe())
		return &prov, nil
	default:
		return nil, errors.Errorf("unknown task pipe " + providerStr)
	}
}
