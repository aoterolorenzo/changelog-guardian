package services

import (
	"bufio"
	"bytes"
	"gitlab.com/aoterocom/changelog-guardian/application/helpers"
	models2 "gitlab.com/aoterocom/changelog-guardian/application/models"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const taskRegexp = `- \[(?P<taskName>[^]]+)?]\((?P<taskHref>[^)]+)?\)\s?(?P<taskTitle>[^@]+)?\s?\(?(@(?P<taskAuthor>[^]]+)?]\((?P<taskAuthorHref>[^)]+)\)\))?`
const releaseRegexp = `## \[(?P<releaseVersion>[^\]]+)]( - (?P<releaseDate>[0-9]{4}-[0-9]{2}-[0-9]{2}))?(?P<releaseYanked> \[YANKED])?`
const releaseLinkRegexp = `\[VERSION]: (?P<releaseLink>.*)`
const categoryRegexp = `### (?P<category>.*)`

func ParseChangelog(pathToChangelog string) (*models2.Changelog, error) {
	c := models2.NewEmptyChangelog()
	changelogReader, err := os.Open(pathToChangelog)
	if err != nil {
		return nil, err
	}

	var currentCategory models2.Category

	fullChangelog, err := ioutil.ReadFile(pathToChangelog)
	if err != nil {
		return nil, err
	}

	changelogReader, err = os.Open(pathToChangelog)
	fscanner := bufio.NewScanner(changelogReader)
	for fscanner.Scan() {
		line := fscanner.Text()
		if bytes.HasPrefix([]byte(line), []byte("## ")) {
			release := models2.NewEmptyRelease()
			release = ParseRelease(line, string(fullChangelog))
			helpers.ReverseAny(c.Releases)
			c.Releases = append(c.Releases, *release)
			helpers.ReverseAny(c.Releases)
		}

		if bytes.HasPrefix([]byte(line), []byte("### ")) {
			currentCategory = ParseCategory(line)
		}

		if bytes.HasPrefix([]byte(line), []byte("- ")) {
			task := models2.NewEmptyTask()
			task = ParseTask(line)
			c.Releases[0].Sections[currentCategory] = append(c.Releases[0].Sections[currentCategory], *task)
		}
	}
	helpers.ReverseAny(c.Releases)
	return c, nil
}

func ParseTask(line string) *models2.Task {

	t := models2.NewEmptyTask()
	r := regexp.MustCompile(taskRegexp)
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
	return t
}

// Map a task string to a Task struct
func ParseRelease(line string, fullChangelog string) *models2.Release {

	r := models2.NewEmptyRelease()
	rg := regexp.MustCompile(releaseRegexp)
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
	rg = regexp.MustCompile(strings.Replace(releaseLinkRegexp, "VERSION", escapedVersion, 1))
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

func ParseCategory(line string) models2.Category {
	r := regexp.MustCompile(categoryRegexp)
	match := r.FindStringSubmatch(line)
	paramsMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return models2.Category(paramsMap["category"])
}
