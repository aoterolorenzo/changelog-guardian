package templates

import (
	"bufio"
	"bytes"
	"fmt"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/helpers"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
)

const markdownTaskRegexp = `- \[(?P<taskName>[^]]+)?]\((?P<taskHref>[^)]+)?\)\s?(?P<taskTitle>[^@]+)?\s?\(?(@(?P<taskAuthor>[^]]+)?]\((?P<taskAuthorHref>[^)]+)\)\))?`
const markdownReleaseRegexp = `## \[(?P<releaseVersion>[^\]]+)]( - (?P<releaseDate>[0-9]{4}-[0-9]{2}-[0-9]{2}))?(?P<releaseYanked> \[YANKED])?`
const markdownReleaseLinkRegexp = `\[VERSION]: (?P<releaseLink>.*)`
const markdownCategoryRegexp = `### (?P<category>.*)`
const markdownChangelogHeader = "# Changelog\n\nAll notable changes to this project will be documented in this file." +
	"\n\nThe format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)," +
	"\nand this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)."

type MarkDownChangelogService struct {
	AbstractChangelog
}

func NewMarkDownChangelogService() *MarkDownChangelogService {
	return &MarkDownChangelogService{AbstractChangelog: AbstractChangelog{}}
}

func (c *MarkDownChangelogService) Parse(pathToChangelog string) (*models.Changelog, error) {
	changelog := models.NewChangelog()
	changelogReader, err := os.Open(pathToChangelog)
	if err != nil {
		return nil, err
	}

	var currentCategory models.Category

	fullChangelog, err := ioutil.ReadFile(pathToChangelog)
	if err != nil {
		return nil, err
	}

	changelogReader, err = os.Open(pathToChangelog)
	fscanner := bufio.NewScanner(changelogReader)
	for fscanner.Scan() {
		line := fscanner.Text()
		if bytes.HasPrefix([]byte(line), []byte("## ")) {
			release := models.NewEmptyRelease()
			release = c.parseRelease(line, string(fullChangelog))
			helpers.ReverseAny(changelog.Releases)
			changelog.Releases = append(changelog.Releases, *release)
			helpers.ReverseAny(changelog.Releases)
		}

		if bytes.HasPrefix([]byte(line), []byte("### ")) {
			currentCategory = c.parseCategory(line)
		}

		if bytes.HasPrefix([]byte(line), []byte("- ")) {
			task := models.NewEmptyTask()
			*task = c.parseTask(line, string(fullChangelog))
			changelog.Releases[0].Sections[currentCategory] = append(changelog.Releases[0].Sections[currentCategory], *task)
		}
	}
	helpers.ReverseAny(changelog.Releases)
	return changelog, nil
}

func (c *MarkDownChangelogService) String(changelog models.Changelog) string {

	changelogStr := markdownChangelogHeader
	changelogStr += "\n"
	changelogStr += c.NudeChangelogString(changelog)
	return changelogStr

}

func (c *MarkDownChangelogService) NudeChangelogString(changelog models.Changelog) string {
	var changelogStr string
	for _, release := range changelog.Releases {
		changelogStr += c.ReleaseToString(release)
	}

	if len(changelog.Releases) > 0 {
		changelogStr += "\n"
		for i, release := range changelog.Releases {
			changelogStr += "[" + release.Version + "]: " + release.Link
			if i != len(changelog.Releases)-1 {
				changelogStr += fmt.Sprintln()
			}
		}
		changelogStr += "\n"
	}

	return changelogStr
}

func (c *MarkDownChangelogService) ReleaseToString(r models.Release) string {
	var yankedStr string

	if r.Yanked {
		yankedStr = " [YANKED]"
	}
	var dateStr string
	if r.Date != "" {
		dateStr = " - " + r.Date
	}

	if strings.ToUpper(r.Version) == "Unreleased" {
		dateStr = ""
	}

	releaseStr := fmt.Sprintf("\n## [%s]%s%s", r.Version, dateStr, yankedStr)
	releaseStr += "\n"

	keys := make([]string, 0, len(r.Sections))
	for k := range r.Sections {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	for _, key := range keys {
		releaseStr += fmt.Sprintf("\n### %s\n\n", key)
		for _, task := range r.Sections[models.Category(key)] {
			releaseStr += fmt.Sprintf(c.TaskToString(task, ""))
			releaseStr += fmt.Sprintln()
		}
	}
	return releaseStr
}

func (c *MarkDownChangelogService) TaskToString(t models.Task, category models.Category) string {
	var authorString string
	if t.Author != "" {
		authorString = fmt.Sprintf(" ([@%s](%s))", t.Author, t.AuthorHref)
	}
	return fmt.Sprintf("- [%s](%s) %s%s", t.ID, t.Href, t.Title, authorString)
}

func (c *MarkDownChangelogService) SaveChangelog(changelog models.Changelog, filePath string) error {
	f, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(c.String(changelog))

	if err != nil {
		return err
	}

	return nil

}

func (c *MarkDownChangelogService) parseTask(line string, fullChangelog string) models.Task {

	t := models.NewEmptyTask()
	r := regexp.MustCompile(markdownTaskRegexp)
	match := r.FindStringSubmatch(line)
	paramsMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	t.ID = paramsMap["taskName"]
	t.Name = paramsMap["taskName"]
	t.Href = paramsMap["taskHref"]
	t.Title = strings.TrimSpace(strings.ReplaceAll(paramsMap["taskTitle"], "([", ""))
	t.Author = paramsMap["taskAuthor"]
	t.AuthorHref = paramsMap["taskAuthorHref"]
	return *t
}

func (c *MarkDownChangelogService) parseRelease(line string, fullChangelog string) *models.Release {

	r := models.NewEmptyRelease()
	rg := regexp.MustCompile(markdownReleaseRegexp)
	match := rg.FindStringSubmatch(line)
	paramsMap := make(map[string]string)
	for i, name := range rg.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	r.Version = paramsMap["releaseVersion"]
	r.Date = paramsMap["releaseDate"]
	if paramsMap["releaseYanked"] != "" {
		r.Yanked = true
	}

	m := regexp.MustCompile("(\\\\|\\^|\\.|\\||\\?|\\*|\\+|\\{|\\}|\\(|\\)|\\[|\\])")
	escapedVersion := m.ReplaceAllString(r.Version, "\\$1")
	rg = regexp.MustCompile(strings.Replace(markdownReleaseLinkRegexp, "VERSION", escapedVersion, 1))
	match = rg.FindStringSubmatch(fullChangelog)
	paramsMap = make(map[string]string)
	for i, name := range rg.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	r.Link = paramsMap["releaseLink"]
	return r
}

func (c *MarkDownChangelogService) parseCategory(line string) models.Category {
	r := regexp.MustCompile(markdownCategoryRegexp)
	match := r.FindStringSubmatch(line)
	paramsMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return models.Category(paramsMap["category"])
}
