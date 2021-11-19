package models

type Task struct {
	Name       string
	Link       string
	Author     string
	AuthorLink string
	Labels     []string
	Files      []string
}

func NewTask(name string, link string, author string, authorLink string, labels []string, files []string) *Task {
	return &Task{Name: name, Link: link, Author: author, AuthorLink: authorLink, Labels: labels, Files: files}
}
