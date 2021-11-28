package settings

import (
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
	ChangelogPath    string `yaml:"changelogPath"`
	ReleaseNotesPath string `yaml:"releaseNotesPath"`
	MainBranch       string `yaml:"mainBranch"`
	DevelopBranch    string `yaml:"defaultBranch"`
	CGConfigPath     string `yaml:"cgConfigPath"`
	Providers        struct {
		Gitlab struct {
			Labels map[models.Category]string `yaml:"labels"`
		} `yaml:"gitlab"`
	} `yaml:"providers"`
	ReleaseProvider string   `yaml:"releaseProvider"`
	TasksProvider   string   `yaml:"tasksProvider"`
	ReleasePipes    []string `yaml:"releasePipes"`
	TasksPipes      []string `yaml:"tasksPipes"`
	TasksPipesCfg   struct {
		ConventionalCommits struct {
			Categories map[models.Category]string `yaml:"categories"`
		} `yaml:"conventional_commits"`
	} `yaml:"tasksPipesCfg"`
	InitialVersion string `yaml:"initialVersion"`
	Style          string `yaml:"style"`
	StylesConfig   struct {
		StylishMarkdown struct {
			Categories map[models.Category][]string `yaml:"categories"`
		} `yaml:"stylish_markdown"`
	} `yaml:"stylesCfg"`
}

func init() {
	Log.Debugf("Generating internal settings...\n")
	if err := extractSettings(settingsFile); err != nil {
		log.Panicln(err)
	}

	Log.Debugf("Retrieving settings from %s...\n", Settings.CGConfigPath)
	yamlFile, err := ioutil.ReadFile(Settings.CGConfigPath)
	if err != nil {
		Log.Debugf("File %s not available. Skilling\n", Settings.CGConfigPath)
	} else {
		if err := extractSettings(string(yamlFile)); err != nil {
			log.Panicln(err)
		}
	}
	Log.Debugf("Settings successfully generated\n")
}

func (g *GlobalSettings) RetrieveSettingsFromFile(file string) {
	settingsFile, err := ioutil.ReadFile(file)

	if err := extractSettings(string(settingsFile)); err != nil {
		log.Panicln(err)
	}

	yamlFile, err := ioutil.ReadFile(Settings.CGConfigPath)
	if err != nil {
	} else {
		if err := extractSettings(string(yamlFile)); err != nil {
			log.Panicln(err)
		}
	}
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
