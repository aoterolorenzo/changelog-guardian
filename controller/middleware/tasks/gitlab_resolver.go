package middleware

import (
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"strings"
)

type GitlabResolverTaskFilter struct {
}

func NewGitlabResolverTaskFilter() *GitlabResolverTaskFilter {
	return &GitlabResolverTaskFilter{}
}

func (tf *GitlabResolverTaskFilter) Filter(task *infra.Task) (*infra.Task, bool, error) {

	if strings.HasPrefix(task.Title, "Resolve \"") && strings.HasSuffix(task.Title, "\"") {
		task.Title = strings.Replace(task.Title, "Resolve \"", "", 1)
		task.Title = strings.TrimSuffix(task.Title, "\"")
		return task, true, nil
	}

	return task, false, nil
}
