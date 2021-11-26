package interfaces

import "gitlab.com/aoterocom/changelog-guardian/application/models"

type ChangelogService interface {
	Parse(pathToChangelog string) (*models.Changelog, error)
	String(changelog models.Changelog) string
	NudeChangelogString(changelog models.Changelog) string
	ReleaseToString(r models.Release) string
	TaskToString(t models.Task, category models.Category) string
	SaveChangelog(changelog models.Changelog, filePath string) error
}
