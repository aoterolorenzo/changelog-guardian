package models

import "time"

type Hash string

type Release struct {
	Name string
	Hash Hash
	Time time.Time
	Link string
}

func NewRelease(name string, hash Hash, time time.Time, link string) *Release {
	return &Release{
		Name: name,
		Hash: hash,
		Time: time,
		Link: link,
	}
}
