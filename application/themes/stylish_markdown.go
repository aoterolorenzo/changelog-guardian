package themes

import (
	"bufio"
	"bytes"
	"fmt"
	"gitlab.com/aoterocom/changelog-guardian/application/helpers"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	settings "gitlab.com/aoterocom/changelog-guardian/config"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
)

const stylishMarkdownTaskRegexp = `-\s[^\s]+\s\[(?P<taskName>[^]]+)?]\((?P<taskHref>[^)]+)?\)\s?(?P<taskTitle>[^@]+)?\s?\(?(@(?P<taskAuthor>[^]]+)?]\((?P<taskAuthorHref>[^)]+)\)\))?`
const stylishMarkdownReleaseRegexp = `## \[\!\[(?P<releaseVersion>[^\]]+)](\s?\!\[)?(?P<releaseDate>[0-9]{4}-[0-9]{2}-[0-9]{2})?]?]\(?(?P<releaseLink>[^)]+)?\)?(?P<releaseYanked> \!\[YANKED])?`
const stylishMarkdownCategoryRegexp = `### \!\[(?P<category>.*)\]`
const stylishMarkdownChangelogHeader = "# Changelog\n\nAll notable changes to this project will be documented in this file." +
	"\n\nThe format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)," +
	"\nand this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)."

type StylishMarkDownChangelogService struct {
	AbstractChangelog
}

func NewStylishMarkDownChangelogService() *StylishMarkDownChangelogService {
	return &StylishMarkDownChangelogService{AbstractChangelog: AbstractChangelog{}}
}

func (c *StylishMarkDownChangelogService) Parse(pathToChangelog string) (*models.Changelog, error) {
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

func (c *StylishMarkDownChangelogService) String(changelog models.Changelog) string {

	changelogStr := stylishMarkdownChangelogHeader
	changelogStr += "\n"
	changelogStr += c.NudeChangelogString(changelog)
	return changelogStr
}

func (c *StylishMarkDownChangelogService) NudeChangelogString(changelog models.Changelog) string {
	var changelogStr string
	for _, release := range changelog.Releases {
		changelogStr += c.ReleaseToString(release)
	}

	for _, release := range changelog.Releases {
		dashedVersion := strings.ReplaceAll(release.Version, "", "")
		urlEscapedDashedVersion := url.QueryEscape(dashedVersion)
		var presentStr = "Release"
		if release.Version == "Unreleased" {
			presentStr = ""
		}
		changelogStr += fmt.Sprintf("\n[%s]: https://img.shields.io/badge/%s-%s-blueviolet?&style=for-the-badge", release.Version, presentStr, urlEscapedDashedVersion)
		if release.Date != "" {
			changelogStr += fmt.Sprintf("\n[%s]: https://img.shields.io/badge/-%s-white?&style=for-the-badge", release.Date, strings.ReplaceAll(release.Date, "-", "--"))
		}
	}

	keys := make([]string, 0, len(settings.Settings.StylesConfig.StylishMarkdown.Categories))
	for k := range settings.Settings.StylesConfig.StylishMarkdown.Categories {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	for _, val := range keys {
		changelogStr += fmt.Sprintf("\n[%s]: https://img.shields.io/badge/-%s-%s.svg?&style=flat-square", val, url.QueryEscape(strings.ToUpper(string(val))), settings.Settings.StylesConfig.StylishMarkdown.Categories[models.Category(val)][0])
	}
	changelogStr += fmt.Sprintf("\n[YANKED]: https://img.shields.io/badge/-YANKED-blueviolet.svg?&style=for-the-badge")
	changelogStr += "\n"

	return changelogStr
}

func (c *StylishMarkDownChangelogService) ReleaseToString(r models.Release) string {
	var yankedStr string

	if r.Yanked {
		yankedStr = " ![YANKED]"
	}
	var dateStr string
	if r.Date != "" {
		dateStr = "![" + r.Date + "]"
	}

	if strings.ToUpper(r.Version) == "Unreleased" {
		dateStr = ""
	}

	releaseStr := fmt.Sprintf("\n## [![%s]%s](%s)%s", r.Version, dateStr, r.Link, yankedStr)
	releaseStr += "\n"

	keys := make([]string, 0, len(r.Sections))
	for k := range r.Sections {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	for _, key := range keys {
		releaseStr += fmt.Sprintf("\n### ![%s]\n\n", key)
		for _, task := range r.Sections[models.Category(key)] {
			releaseStr += fmt.Sprintf(c.TaskToString(task, models.Category(key)))
			releaseStr += fmt.Sprintln()
		}
	}
	return releaseStr
}

func (c *StylishMarkDownChangelogService) TaskToString(t models.Task, category models.Category) string {
	var authorString string
	if t.Author != "" {

		authorString = fmt.Sprintf(" ([@%s](%s))", t.Author, t.AuthorHref)
	}

	var categoryEmoji = settings.Settings.StylesConfig.StylishMarkdown.Categories[category][1]

	return fmt.Sprintf("- %s [%s](%s) %s%s", categoryEmoji, t.ID, t.Href, t.Title, authorString)
}

func (c *StylishMarkDownChangelogService) SaveChangelog(changelog models.Changelog, filePath string) error {
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

func (c *StylishMarkDownChangelogService) parseTask(line string, fullChangelog string) models.Task {

	t := models.NewEmptyTask()
	r := regexp.MustCompile(stylishMarkdownTaskRegexp)
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

func (c *StylishMarkDownChangelogService) parseRelease(line string, fullChangelog string) *models.Release {

	r := models.NewEmptyRelease()
	rg := regexp.MustCompile(stylishMarkdownReleaseRegexp)
	match := rg.FindStringSubmatch(line)
	paramsMap := make(map[string]string)
	for i, name := range rg.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	r.Version = paramsMap["releaseVersion"]
	r.Date = paramsMap["releaseDate"]
	r.Link = paramsMap["releaseLink"]
	if paramsMap["releaseYanked"] != "" {
		r.Yanked = true
	}
	return r
}

func (c *StylishMarkDownChangelogService) parseCategory(line string) models.Category {
	r := regexp.MustCompile(stylishMarkdownCategoryRegexp)
	match := r.FindStringSubmatch(line)
	paramsMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return models.Category(paramsMap["category"])
}
