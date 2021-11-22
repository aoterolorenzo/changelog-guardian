package interfaces

import (
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
)

type ReleaseFilter interface {
	Filter(release *infra.Release) (*infra.Release, bool, error)
}

type TaskFilter interface {
	Filter(task *infra.Task) (*infra.Task, bool, error)
}
