package providers

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	settings "gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/helpers"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"gopkg.in/andygrunwald/go-jira.v1"
	"os"
)

type JiraController struct {
	user    string
	token   string
	baseUrl string
}

func NewJiraController() *JiraController {
	return &JiraController{
		user:    os.Getenv("JIRA_USER"),
		token:   os.Getenv("JIRA_API_TOKEN"),
		baseUrl: settings.Settings.TasksPipesCfg.Jira.BaseUrl,
	}
}

func (jc JiraController) GetTask(taskId string) (*infra.Task, error) {

	var jiraURL = settings.Settings.TasksPipesCfg.Jira.BaseUrl

	tp := jira.BasicAuthTransport{
		Username: jc.user,
		Password: jc.token,
	}

	client, err := jira.NewClient(tp.Client(), jiraURL)
	issue, res, err := client.Issue.Get(taskId, nil)
	if err != nil {
		return &infra.Task{}, err
	}

	if res.StatusCode != 200 {
		return &infra.Task{}, errors.New(fmt.Sprintf("error retrieving the task. status code %d", res.StatusCode))
	}

	issue, res, err = client.Issue.Get(taskId, nil)
	if err != nil {
		return &infra.Task{}, err
	}

	return &infra.Task{
		ID:         issue.Key,
		Name:       issue.Fields.Summary,
		Link:       jc.baseUrl + "/browse/" + issue.Key,
		Title:      issue.Key,
		Author:     issue.Fields.Creator.DisplayName,
		AuthorLink: jc.baseUrl + "/jira/people/" + issue.Fields.Creator.AccountID,
		Category:   jc.inferCategory(issue.Fields.Labels, issue.Fields.Type),
	}, nil
}

func (jc JiraController) inferCategory(labels []string, taskType jira.IssueType) models.Category {

	// If ticket has specific labels
	for _, val := range settings.Settings.TasksPipesCfg.Jira.Labels {
		if helpers.SliceContainsString(labels, val) {
			return models.Category(val)
		}
	}

	// If ticket is type bug or any custom type with the bug icon
	if taskType.ID == "10004" || taskType.IconURL ==
		"https://aoterocom.atlassian.net/rest/api/2/universal_avatar/view/type/issuetype/avatar/10303?size=medium" {
		return models.FIXED
	}

	// Otherwise, it's just a task
	return models.ADDED
}
