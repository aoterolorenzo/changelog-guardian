package services

import (
	"bufio"
	"bytes"
	"gitlab.com/aoterocom/changelog-guardian/models"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
)

const taskRegexp = `- \[(?P<taskName>[^]]+)?]\((?P<taskHref>[^)]+)?\)\s?(?P<taskTitle>[^@]+)?\s?\(?(@(?P<taskAuthor>[^]]+)?]\((?P<taskAuthorHref>[^)]+)\)\))?`
const releaseRegexp = `## \[(?P<releaseVersion>[^\]]+)]( - (?P<releaseDate>[0-9]{4}-[0-9]{2}-[0-9]{2}))?(?P<releaseYanked> \[YANKED])?`
const releaseLinkRegexp = `\[VERSION]: (?P<releaseLink>.*)`
const categoryRegexp = `### (?P<category>.*)`

func ParseChangelog(pathToChangelog string) *models.Changelog {
	c := models.NewEmptyChangelog()
	changelogReader, err := os.Open(pathToChangelog)
	if err != nil {
		panic(err)
	}

	var currentCategory models.Category

	fullChangelog, err := ioutil.ReadFile(pathToChangelog)
	if err != nil {
		panic(err)
	}

	changelogReader, err = os.Open(pathToChangelog)
	fscanner := bufio.NewScanner(changelogReader)
	for fscanner.Scan() {
		line := fscanner.Text()
		if bytes.HasPrefix([]byte(line), []byte("## ")) {
			release := models.NewEmptyRelease()
			release = ParseRelease(line, string(fullChangelog))
			reverseAny(c.Releases)
			c.Releases = append(c.Releases, *release)
			reverseAny(c.Releases)
		}

		if bytes.HasPrefix([]byte(line), []byte("### ")) {
			currentCategory = ParseCategory(line)
		}

		if bytes.HasPrefix([]byte(line), []byte("- ")) {
			task := models.NewEmptyTask()
			task = ParseTask(line)
			c.Releases[0].Sections[currentCategory] = append(c.Releases[0].Sections[currentCategory], *task)
		}
	}
	reverseAny(c.Releases)
	return c
}

func ParseTask(line string) *models.Task {

	t := models.NewEmptyTask()
	r := regexp.MustCompile(taskRegexp)
	match := r.FindStringSubmatch(line)
	paramsMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	t.Name = paramsMap["taskName"]
	t.Href = paramsMap["taskHref"]
	t.Title = strings.TrimSpace(strings.ReplaceAll(paramsMap["taskTitle"], "([", ""))
	t.Author = paramsMap["taskAuthor"]
	t.AuthorHref = paramsMap["taskAuthorHref"]
	return t
}

// Map a task string to a Task struct
func ParseRelease(line string, fullChangelog string) *models.Release {

	r := models.NewEmptyRelease()
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

	rg = regexp.MustCompile(strings.Replace(releaseLinkRegexp, "VERSION", r.Version, 1))
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

func ParseCategory(line string) models.Category {
	r := regexp.MustCompile(categoryRegexp)
	match := r.FindStringSubmatch(line)
	paramsMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return models.Category(paramsMap["category"])
}

func reverseAny(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
