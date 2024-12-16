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
		Git struct {
			GitRoot string `yaml:"gitRoot"`
		} `yaml:"git"`
		Gitlab struct {
			Labels  map[models.Category]string `yaml:"labels"`
			GitRoot string                     `yaml:"gitRoot"`
		} `yaml:"gitlab"`
		Github struct {
			Labels  map[models.Category]string `yaml:"labels"`
			GitRoot string                     `yaml:"gitRoot"`
		} `yaml:"github"`
		GithubPRs struct {
			GHReleaseSearch string `yaml:"ghReleaseSearch"`
			VersionRegex    string `yaml:"versionRegex"`
			TargetBranch    string `yaml:"targetBranch"`
			GitRoot         string `yaml:"gitRoot"`
		} `yaml:"githubPRs"`
	} `yaml:"providers"`
	ReleaseProvider   string   `yaml:"releaseProvider"`
	TasksProvider     string   `yaml:"tasksProvider"`
	CategoryFromPipes bool     `yaml:"categoryFromPipes"`
	ReleasePipes      []string `yaml:"releasePipes"`
	TasksPipes        []string `yaml:"tasksPipes"`
	TasksPipesCfg     struct {
		ConventionalCommits struct {
			Categories map[models.Category]string `yaml:"categories"`
		} `yaml:"conventional_commits"`
		Jira struct {
			BaseUrl string                     `yaml:"baseUrl"`
			REGEX   string                     `yaml:"regex"`
			Labels  map[models.Category]string `yaml:"labels"`
		} `yaml:"jira"`
		InclusionsExclusions struct {
			Labels struct {
				Inclusions []string `yaml:"included"`
				Exclusions []string `yaml:"excluded"`
			} `yaml:"labels"`
			Regexps struct {
				Inclusions []string `yaml:"included"`
				Exclusions []string `yaml:"excluded"`
			} `yaml:"regexps"`
			Paths struct {
				Inclusions []string `yaml:"included"`
				Exclusions []string `yaml:"excluded"`
			} `yaml:"paths"`
		} `yaml:"inclusions_exclusions"`
	} `yaml:"tasksPipesCfg"`
	InitialVersion  string `yaml:"initialVersion"`
	Template        string `yaml:"template"`
	TemplatesConfig struct {
		StylishMarkdown struct {
			Categories map[models.Category][]string `yaml:"categories"`
		} `yaml:"stylish_markdown"`
	} `yaml:"templatesCfg"`
}

func init() {
	Log.Debugf("Generating internal settings...\n")
	if err := extractSettings(settingsFile); err != nil {
		log.Panicln(err)
	}

	Log.Debugf("Retrieving settings from %s...\n", Settings.CGConfigPath)
	yamlFile, err := ioutil.ReadFile(Settings.CGConfigPath)
	if err != nil {
		Log.WithError(err).Debugf("File %s not available. Skipping\n", Settings.CGConfigPath)
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
