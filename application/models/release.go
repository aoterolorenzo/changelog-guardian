package models

import (
	"fmt"
	"sort"
	"strings"
)

type Release struct {
	Version  string
	Date     string
	Link     string
	Yanked   bool
	Sections map[Category][]Task
}

func NewRelease(version string, date string, link string, yanked bool, sections map[Category][]Task) *Release {
	if sections == nil {
		sections = make(map[Category][]Task)
	}
	return &Release{
		Version:  version,
		Date:     date,
		Link:     link,
		Yanked:   yanked,
		Sections: sections}
}

func NewEmptyRelease() *Release {
	return &Release{
		Sections: make(map[Category][]Task),
	}
}

// Stringify the release
func (r *Release) String() string {
	var yankedStr string

	if r.Yanked {
		yankedStr = " [YANKED]"
	}
	var dateStr string
	if r.Date != "" {
		dateStr = " - " + r.Date
	}

	if strings.ToUpper(r.Version) == "UNRELEASED" {
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
		for _, task := range r.Sections[Category(key)] {
			releaseStr += fmt.Sprintf(task.String())
			releaseStr += fmt.Sprintln()
		}
	}
	return releaseStr
}
