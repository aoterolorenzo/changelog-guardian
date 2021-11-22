package settings

import (
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gopkg.in/yaml.v3"
	"log"
)

import _ "embed"

//go:embed settings.yml
var settingsFile string

var Settings = &GlobalSettings{}

type GlobalSettings struct {
	ChangelogPath string `yaml:"changelogPath"`
	MainBranch    string `yaml:"mainBranch"`
	DevelopBranch string `yaml:"developBranch"`
	Providers     struct {
		Gitlab struct {
			Labels map[models.Category]string `yaml:"labels"`
		} `yaml:"gitlab"`
	} `yaml:"providers"`
}

func init() {
	if err := extractSettings(settingsFile); err != nil {
		log.Panicln(err)
	}
}

func extractSettings(file string) error {
	err := yaml.Unmarshal([]byte(file), Settings)
	return err
}
