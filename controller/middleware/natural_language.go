package middleware

import (
	application "gitlab.com/aoterocom/changelog-guardian/application/models"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
)

type NaturalLanguageMiddleware struct {
}

func (nlm *NaturalLanguageMiddleware) Filter(task infra.Task) (*application.Task, error) {
	return nil, nil
}
