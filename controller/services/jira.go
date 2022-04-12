package services

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	settings "gitlab.com/aoterocom/changelog-guardian/config"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"gopkg.in/andygrunwald/go-jira.v1"
	"os"
)

type JiraService struct {
	user    string
	token   string
	baseUrl string
}

func NewJiraService() *JiraService {
	return &JiraService{
		user:    os.Getenv("JIRA_USER"),
		token:   os.Getenv("JIRA_API_TOKEN"),
		baseUrl: settings.Settings.TasksPipesCfg.Jira.BaseUrl,
	}
}

func (jc JiraService) GetTask(taskId string) (infra.Task, error) {

	var jiraURL = settings.Settings.TasksPipesCfg.Jira.BaseUrl

	tp := jira.BasicAuthTransport{
		Username: jc.user,
		Password: jc.token,
	}

	client, err := jira.NewClient(tp.Client(), jiraURL)
	issue, res, err := client.Issue.Get(taskId, nil)
	if err != nil {
		return infra.Task{}, err
	}

	if res.StatusCode != 200 {
		return infra.Task{}, errors.New(fmt.Sprintf("error retrieving the task. status code %d", res.StatusCode))
	}

	issue, res, err = client.Issue.Get(taskId, nil)

	if err != nil {
		return infra.Task{}, err
	}
	return infra.Task{
		ID:         issue.Key,
		Name:       issue.Fields.Summary,
		Link:       jc.baseUrl + "/browse/" + issue.Key,
		Title:      issue.Key,
		Author:     issue.Fields.Creator.DisplayName,
		AuthorLink: jc.baseUrl + "/jira/people/" + issue.Fields.Creator.AccountID,
		Category:   inferCategory(issue.Fields.Labels, issue.Fields.Type.ID),
	}, nil
}

func inferCategory(labels []string, taskType string) models.Category {
	// TODO: implement category inference
	return models.ADDED
}
