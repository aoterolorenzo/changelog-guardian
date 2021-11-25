package interfaces

import "gitlab.com/aoterocom/changelog-guardian/application/models"

type ChangelogService interface {
	Parse(pathToChangelog string) (*models.Changelog, error)
	String(changelog models.Changelog) string
	SaveChangelog(changelog models.Changelog, filePath string) error
}
