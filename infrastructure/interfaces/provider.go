package interfaces

import (
	application "gitlab.com/aoterocom/changelog-guardian/application/models"
	infrastructure "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"time"
)

type Provider interface {
	GetReleases(from *time.Time, to *time.Time, repo *string) (*[]infrastructure.Release, error)
	GetTasks(from *time.Time, to *time.Time, repo *string, targetBranch string) (*[]infrastructure.Task, error)
	DefineCategory(task infrastructure.Task) application.Category
	ReleaseURL(repo string, from *string, to string) string
}
