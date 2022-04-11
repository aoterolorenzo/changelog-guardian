package services

import (
	"fmt"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	settings "gitlab.com/aoterocom/changelog-guardian/config"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"golang.org/x/oauth2"
	jiraOauth "golang.org/x/oauth2/jira"
	jira "gopkg.in/andygrunwald/go-jira.v1"
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
		token:   os.Getenv("JIRA_TOKEN"),
		baseUrl: settings.Settings.TasksPipesCfg.Jira.BaseUrl,
	}
}

func (jc JiraService) GetTask(taskId string) (infra.Task, error) {

	var (
		jiraURL = settings.Settings.TasksPipesCfg.Jira.BaseUrl
	)

	var conf = jiraOauth.Config{
		BaseURL: jiraURL,
		Subject: "urn:atlassian:connect:useraccountid:623df10b8d8b9c0068b9b21e",
		Config: oauth2.Config{
			ClientID:     "qYLJ1eAUnchVsYjUS5S5L8CVRpehCHUg",
			ClientSecret: "XunzfAxgg6nW52BFoJOlDtgrHIAIXvnIgQ_ibwakf1S2SV4F9tu0v3I4v8hPAsCT",
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://auth.atlassian.com/authorize",
				TokenURL:  "https://auth.atlassian.com/oauth/token",
				AuthStyle: 0,
			},
			RedirectURL: "https://aoterocom.atlassian.com/",
			Scopes:      []string{"read", "write"},
		},
	}

	fmt.Println("get task " + taskId)

	jiraClient, err := jira.NewClient(conf.Client(nil), jc.baseUrl+"jira/")
	if err != nil {
		return infra.Task{}, err
	}

	issue, res, err := jiraClient.Issue.Get(taskId, nil)

	fmt.Println(res)
	fmt.Println(issue)
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
