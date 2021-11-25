package services

import (
	"github.com/pkg/errors"
	"gitlab.com/aoterocom/changelog-guardian/controller/interfaces"
	releasePipes "gitlab.com/aoterocom/changelog-guardian/controller/pipes/release"
	taskPipes "gitlab.com/aoterocom/changelog-guardian/controller/pipes/tasks"
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

func TaskPipeSelector(providerStr string) (*interfaces.TaskPipe, error) {
	switch providerStr {
	case "gitlab_resolver":
		prov := interfaces.TaskPipe(taskPipes.NewGitlabResolverTaskPipe())
		return &prov, nil
	case "natural_language":
		prov := interfaces.TaskPipe(taskPipes.NewNaturalLanguageTaskPipe())
		return &prov, nil
	default:
		return nil, errors.Errorf("unknown task pipe " + providerStr)
	}
}
