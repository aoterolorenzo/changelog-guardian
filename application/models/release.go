package models

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
