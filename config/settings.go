package settings

import (
	"fmt"
	"github.com/imdario/mergo"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

import _ "embed"

//go:embed settings.yml
var settingsFile string

var Settings = &GlobalSettings{}

type GlobalSettings struct {
	ChangelogPath string `yaml:"changelogPath"`
	MainBranch    string `yaml:"mainBranch"`
	DevelopBranch string `yaml:"defaultBranch"`
	CGConfigPath  string `yaml:"cgConfigPath"`
	Providers     struct {
		Gitlab struct {
			Labels map[models.Category]string `yaml:"labels"`
		} `yaml:"gitlab"`
	} `yaml:"providers"`
	ReleaseProvider string   `yaml:"releaseProvider"`
	TasksProvider   string   `yaml:"tasksProvider"`
	ReleaseFilters  []string `yaml:"releaseFilters"`
	TaskFilters     []string `yaml:"taskFilters"`
}

func init() {
	fmt.Println("Constructing internal settings...")
	if err := extractSettings(settingsFile); err != nil {
		log.Panicln(err)
	}

	fmt.Println("Retrieving settings from " + Settings.CGConfigPath + "...")
	yamlFile, err := ioutil.ReadFile(Settings.CGConfigPath)
	if err != nil {
		fmt.Printf(Settings.CGConfigPath + " not available. Skipping...")
	} else {
		if err := extractSettings(string(yamlFile)); err != nil {
			log.Panicln(err)
		}
	}

	fmt.Println("Settings successfully retrieved.")
}

func extractSettings(content string) error {
	var setts = &GlobalSettings{}
	if err := yaml.Unmarshal([]byte(content), setts); err != nil {
		return err
	}

	if err := mergo.Merge(Settings, setts, mergo.WithOverride); err != nil {
		return err
	}

	return nil
}
