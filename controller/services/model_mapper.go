package services

import (
	application "gitlab.com/aoterocom/changelog-guardian/application/models"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
)

type ModelMapperService struct {
}

func NewModelMapperService() *ModelMapperService {
	return &ModelMapperService{}
}

func (mm *ModelMapperService) InfraReleaseToApplicationModel(release infra.Release) *application.Release {
	return application.NewRelease(release.Name, release.Time.Format("2006-01-02"), release.Link, false, nil)
}

func (mm *ModelMapperService) InfraTaskToApplicationModel(task infra.Task) *application.Task {
	return application.NewTask(task.ID, task.Name, task.Link, task.Title, task.Author, task.AuthorLink, task.Category)
}
