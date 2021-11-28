package selectors

import (
	"github.com/pkg/errors"
	"gitlab.com/aoterocom/changelog-guardian/controller/interfaces"
	releasePipes "gitlab.com/aoterocom/changelog-guardian/controller/pipes/release"
	tasksPipes "gitlab.com/aoterocom/changelog-guardian/controller/pipes/tasks"
)

func ReleasePipeSelector(providerStr string) (*interfaces.ReleasePipe, error) {
	switch providerStr {
	case "semver":
		prov := interfaces.ReleasePipe(releasePipes.NewSemverReleasePipe())
		return &prov, nil
	default:
		return nil, errors.Errorf("unknown release pipe " + providerStr)
	}
}

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
	default:
		return nil, errors.Errorf("unknown task pipe " + providerStr)
	}
}
