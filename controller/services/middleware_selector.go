package services

import (
	"github.com/pkg/errors"
	"gitlab.com/aoterocom/changelog-guardian/controller/interfaces"
	releaseFilters "gitlab.com/aoterocom/changelog-guardian/controller/middleware/release"
	taskFilters "gitlab.com/aoterocom/changelog-guardian/controller/middleware/tasks"
)

func ReleaseFilterSelector(providerStr string) (*interfaces.ReleaseFilter, error) {
	switch providerStr {
	case "semver":
		prov := interfaces.ReleaseFilter(releaseFilters.NewSemverReleaseFilter())
		return &prov, nil
	default:
		return nil, errors.Errorf("unknown release filter " + providerStr)
	}
}

func TaskFilterSelector(providerStr string) (*interfaces.TaskFilter, error) {
	switch providerStr {
	case "gitlab_resolver":
		prov := interfaces.TaskFilter(taskFilters.NewGitlabResolverTaskFilter())
		return &prov, nil
	case "natural_language":
		prov := interfaces.TaskFilter(taskFilters.NewNaturalLanguageTaskFilter())
		return &prov, nil
	default:
		return nil, errors.Errorf("unknown task filter " + providerStr)
	}
}
