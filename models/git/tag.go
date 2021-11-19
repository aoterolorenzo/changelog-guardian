package models

import "time"

type Hash string

type Tag struct {
	Name string
	Hash Hash
	Time time.Time
}

func NewTag(name string, hash Hash, time time.Time) *Tag {
	return &Tag{
		Name: name,
		Hash: hash,
		Time: time,
	}
}
