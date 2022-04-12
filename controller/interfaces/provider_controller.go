package interfaces

import (
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
)

type ProviderController interface {
	GetTask(taskId string) (*infra.Task, error)
}
