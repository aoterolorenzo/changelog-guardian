package models

type Changelog struct {
	Releases []Release
}

func NewChangelog() *Changelog {
	return &Changelog{}
}
