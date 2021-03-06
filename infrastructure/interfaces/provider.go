package interfaces

import (
	application "gitlab.com/aoterocom/changelog-guardian/application/models"
	infrastructure "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"time"
)

type Provider interface {
	GetReleases(from *time.Time, to *time.Time) (*[]infrastructure.Release, error)
	GetTasks(from *time.Time, to *time.Time, targetBranch string) (*[]infrastructure.Task, error)
	GetTask(taskId string) (*infrastructure.Task, error)
	DefineCategory(task infrastructure.Task) application.Category
	ReleaseURL(from *string, to string) (*string, error)
}
