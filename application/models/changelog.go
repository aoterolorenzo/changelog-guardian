package models

import (
	"fmt"
)

const changelogHeader = "# Changelog\n\nAll notable changes to this project will be documented in this file." +
	"\n\nThe format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)," +
	"\nand this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)."

type Changelog struct {
	Releases []Release
}

func NewChangelog(releases []Release) *Changelog {
	return &Changelog{Releases: releases}
}

func NewEmptyChangelog() *Changelog {
	return &Changelog{}
}

func (c *Changelog) String() string {
	changelogStr := changelogHeader
	changelogStr += "\n"
	for _, release := range c.Releases {
		changelogStr += release.String()
	}

	if len(c.Releases) > 0 {
		changelogStr += "\n"
		for i, release := range c.Releases {
			changelogStr += "[" + release.Version + "]: " + release.Link
			if i != len(c.Releases)-1 {
				changelogStr += fmt.Sprintln()
			}
		}
		changelogStr += "\n"
	}
	return changelogStr
}
