package models

import "time"

type Hash string

type Release struct {
	Name string
	Hash Hash
	Time time.Time
}

func NewRelease(name string, hash Hash, time time.Time) *Release {
	return &Release{
		Name: name,
		Hash: hash,
		Time: time,
	}
}
