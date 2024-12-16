package providers

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/andygrunwald/go-jira.v1"

	"gitlab.com/aoterocom/changelog-guardian/application/models"
	settings "gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/helpers"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
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

	var authorName string
	var authorLink string
	author := issue.Fields.Assignee
	if author == nil {
		author = issue.Fields.Creator
	}
	authorLink = jc.baseUrl + "secure/ViewProfile.jspa?name=" + author.Name

	return &infra.Task{
		ID:         issue.Key,
		Name:       issue.Fields.Summary,
		Link:       jc.baseUrl + "/browse/" + issue.Key,
		Title:      issue.Fields.Summary,
		Author:     authorName,
		AuthorLink: authorLink,
		Category:   jc.inferCategory(issue.Fields.Labels, issue.Fields.Type),
	}, nil
}

func (jc JiraController) inferCategory(labels []string, taskType jira.IssueType) models.Category {

	// If ticket has specific labels
	for key, val := range settings.Settings.TasksPipesCfg.Jira.Labels {
		if helpers.SliceContainsString(labels, val) {
			return key
		}
	}

	// If ticket is type bug or any custom type with the bug icon
	if taskType.Name == "Bug" || strings.Contains(taskType.IconURL, "10303") {
		return models.FIXED
	}

	// Otherwise, it's just a task
	return models.ADDED
}
