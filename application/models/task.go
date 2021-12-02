package models

type Task struct {
	ID         string
	Name       string
	Href       string
	Title      string
	Author     string
	AuthorHref string
	Category   Category
}

func NewTask(id string, name string, href string, title string, author string, authorHref string, category Category) *Task {
	return &Task{
		ID:         id,
		Name:       name,
		Href:       href,
		Title:      title,
		Author:     author,
		AuthorHref: authorHref,
		Category:   category,
	}
}

func NewEmptyTask() *Task {
	return &Task{}
}
