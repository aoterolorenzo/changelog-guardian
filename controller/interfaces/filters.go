package interfaces

import (
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
)

type ReleasePipe interface {
	Filter(release *infra.Release) (*infra.Release, bool, error)
}

type TasksPipe interface {
	Filter(task *infra.Task) (*infra.Task, bool, error)
}
