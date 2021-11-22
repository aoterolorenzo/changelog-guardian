package interfaces

import (
	application "gitlab.com/aoterocom/changelog-guardian/application/models"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
)

type Middleware interface {
	Filter(task infra.Task) (*application.Task, error)
}
