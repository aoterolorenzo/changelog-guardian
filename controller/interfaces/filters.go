package interfaces

import (
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
)

type ReleasePipe interface {
	Pipe(release *infra.Release) (*infra.Release, bool, error)
}

type TaskPipe interface {
	Pipe(task *infra.Task) (*infra.Task, bool, error)
}
