package services

import (
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
)

type MockJiraService struct{}

func (jc MockJiraService) GetTask(taskId string) (*infra.Task, error) {
	switch taskId {
	case "TES-1":
		return &infra.Task{Title: "TES1", Category: models.ADDED}, nil

	case "TES-2":
		return &infra.Task{Title: "TES2", Category: models.ADDED}, nil

	}

	return nil, nil
}
