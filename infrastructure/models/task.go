package models

import "gitlab.com/aoterocom/changelog-guardian/application/models"

type Task struct {
	ID         string
	Name       string
	Title      string
	Link       string
	Author     string
	AuthorLink string
	Labels     []string
	Files      []string
	Category   models.Category
}

func NewTask(id string, name string, title string, link string, author string, authorLink string, labels []string, files []string) *Task {
	return &Task{ID: id, Name: name, Title: title, Link: link, Author: author, AuthorLink: authorLink, Labels: labels, Files: files}
}
