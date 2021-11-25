package middleware

import (
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"strings"
)

type NaturalLanguageTaskPipe struct {
}

func NewNaturalLanguageTaskPipe() *NaturalLanguageTaskPipe {
	return &NaturalLanguageTaskPipe{}
}

func (nlm *NaturalLanguageTaskPipe) Pipe(task *infra.Task) (*infra.Task, bool, error) {

	words := [][]string{
		{"Add", "Added"},
		{"Fix", "Fixed"},
		{"Refactor", "Refactorized"},
		{"Change", "Changed"},
		{"Implement", "Implemented"},
		{"Remove", "Removed"},
		{"Document", "Documented"},
		{"Improve", "Improved"},
		{"Finish", "Finished"},
	}

	for _, pair := range words {
		if strings.HasPrefix(task.Title, pair[0]) {
			task.Title = strings.Replace(task.Title, pair[0], pair[1], 1)
			return task, true, nil
		}
	}

	return task, false, nil
}
