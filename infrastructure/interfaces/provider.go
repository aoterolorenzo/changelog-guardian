package interfaces

import (
	models2 "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"time"
)

type Provider interface {
	GetReleases() ([]*models2.Release, error)
	GetTasks(from *time.Time, to *time.Time, targetBranch string) (*[]models2.Task, error)
}
