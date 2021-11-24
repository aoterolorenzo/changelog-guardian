package themes

import (
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"os"
)

type AbstractChangelog struct {
}

func (a *AbstractChangelog) SaveChangelog(changelog models.Changelog, filePath string) error {
	f, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(a.String(changelog))

	if err != nil {
		return err
	}

	return nil

}

func (a *AbstractChangelog) String(changelog models.Changelog) string {
	return ""
}
