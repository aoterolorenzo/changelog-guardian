package providers

import (
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
	application "gitlab.com/aoterocom/changelog-guardian/application/models"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/helpers"
	infrastructure "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type GitlabProvider struct {
	GitToken string
	repo     *string
}

func init() {
	cwd, _ := os.Getwd()
	_ = godotenv.Load(cwd + "/vars.env")
}

func NewGitlabProvider(repo *string) *GitlabProvider {
	return &GitlabProvider{
		GitToken: os.Getenv("GITLAB_TOKEN"),
		repo:     repo,
	}
}

func (gp *GitlabProvider) GetReleases(from *time.Time, to *time.Time) (*[]infrastructure.Release, error) {

	namespacedRepo, err := gp.namespacedRepo()
	if err != nil {
		return nil, err
	}

	gitlabClient, _ := gitlab.NewClient(gp.GitToken)
	project, _, err := gitlabClient.Projects.GetProject(url.QueryEscape(*namespacedRepo), &gitlab.GetProjectOptions{})
	if err != nil {
		id, err := gp.getProjectIdWithDirectRequest(namespacedRepo)
		project = &gitlab.Project{ID: id}
		if err != nil {
			return nil, err
		}
	}

	var releases []*gitlab.Release
	opts := &gitlab.ListReleasesOptions{}

	for {
		// Get the first page with projects.
		releasesRetrieved, resp, err := gitlabClient.Releases.ListReleases(project.ID, opts)
		if err != nil {
			return nil, err
		}

		releases = append(releases, releasesRetrieved...)

		// Exit the loop when we've seen all pages.
		if resp.NextPage == 0 {
			break
		}

		// Update the page number to get the next page.
		opts.Page = resp.NextPage
	}

	var gitReleases []infrastructure.Release
	for _, release := range releases {
		releaseLink := "https://gitlab.com/" + *namespacedRepo + "/-/releases/" + release.Name
		gitReleases = append(gitReleases, *infrastructure.NewRelease(release.Name, infrastructure.Hash(release.Commit.ID), *release.ReleasedAt, releaseLink))
	}
	helpers.ReverseAny(gitReleases)
	return &gitReleases, nil
}

func (gp *GitlabProvider) getProjectIdWithDirectRequest(namespacedRepo *string) (int, error) {

	req, err := http.NewRequest("GET", "https://gitlab.com/api/v4/projects/"+url.QueryEscape(*namespacedRepo), nil)
	if err != nil {
		return 0, err
	}

	client := &http.Client{}
	req.Header.Add("PRIVATE-TOKEN", gp.GitToken)
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New(fmt.Sprintf("response status code: %d", http.StatusOK))
	}

	var gitlabProject struct {
		Id int `json:"id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&gitlabProject)
	if err != nil {
		return 0, err
	}

	return gitlabProject.Id, nil
}

func (gp *GitlabProvider) GetTasks(from *time.Time, to *time.Time, targetBranch string) (*[]infrastructure.Task, error) {

	namespacedRepo, err := gp.namespacedRepo()
	if err != nil {
		return nil, err
	}

	gitlabClient, _ := gitlab.NewClient(gp.GitToken)
	project, _, err := gitlabClient.Projects.GetProject(*namespacedRepo, &gitlab.GetProjectOptions{})
	if err != nil {
		return nil, err
	}

	var gitTasks []infrastructure.Task
	var state = "merged"
	listProjectMergeRequestOptions := &gitlab.ListProjectMergeRequestsOptions{CreatedAfter: from, CreatedBefore: to,
		TargetBranch: &targetBranch, State: &state}
	var mergeRequests []*gitlab.MergeRequest

	for {
		// Get the first page with projects.
		mergeRequestsRetrieved, resp, err := gitlabClient.MergeRequests.ListProjectMergeRequests(project.ID, listProjectMergeRequestOptions)
		if err != nil {
			return nil, err
		}

		mergeRequests = append(mergeRequests, mergeRequestsRetrieved...)

		// Exit the loop when we've seen all pages.
		if resp.NextPage == 0 {
			break
		}

		// Update the page number to get the next page.
		listProjectMergeRequestOptions.Page = resp.NextPage
	}

	for _, mergeRequest := range mergeRequests {
		var labelStrings []string
		labels := mergeRequest.Labels
		for _, label := range labels {
			labelStrings = append(labelStrings, label)
		}

		var fileChanges []string
		for _, change := range mergeRequest.Changes {
			fileChanges = append(fileChanges, change.NewPath)
		}
		gitTask := infrastructure.NewTask("!"+strconv.Itoa(mergeRequest.IID), "!"+strconv.Itoa(mergeRequest.IID), mergeRequest.Title, mergeRequest.WebURL, mergeRequest.Author.Username,
			mergeRequest.Author.WebURL, labelStrings, fileChanges)
		gitTask.Category = gp.DefineCategory(*gitTask)
		gitTasks = append(gitTasks, *gitTask)

	}

	return &gitTasks, nil
}

func (gp *GitlabProvider) DefineCategory(task infrastructure.Task) application.Category {
	var category = application.ADDED

	if strings.HasPrefix(strings.ToLower(task.Title), "revert") {
		return application.REMOVED
	}

	for _, label := range task.Labels {
		switch label {
		case Settings.Providers.Gitlab.Labels[application.ADDED]:
			category = application.ADDED
			break
		case Settings.Providers.Gitlab.Labels[application.FIXED]:
			category = application.FIXED
			break
		case Settings.Providers.Gitlab.Labels[application.REFACTOR]:
			category = application.REFACTOR
			break
		case Settings.Providers.Gitlab.Labels[application.DEPRECATED]:
			category = application.DEPRECATED
			break
		case Settings.Providers.Gitlab.Labels[application.CHANGED]:
			category = application.CHANGED
			break
		case Settings.Providers.Gitlab.Labels[application.DEPENDENCIES]:
			category = application.DEPENDENCIES
			break
		case Settings.Providers.Gitlab.Labels[application.DOCUMENTATION]:
			category = application.DOCUMENTATION
			break
		case Settings.Providers.Gitlab.Labels[application.REMOVED]:
			category = application.REMOVED
			break
		case Settings.Providers.Gitlab.Labels[application.SECURITY]:
			category = application.SECURITY
			break
		case Settings.Providers.Gitlab.Labels[application.BREAKING_CHANGE]:
			category = application.BREAKING_CHANGE
			return category
		default:
			break
		}
	}

	return category
}

func (gp *GitlabProvider) ReleaseURL(from *string, to string) (*string, error) {
	namespacedRepo, err := gp.namespacedRepo()
	if err != nil {
		return nil, err
	}

	if from != nil {
		url := "https://gitlab.com/" + *namespacedRepo + "/-/compare/" + *from + "..." + to
		return &url, nil
	}

	url := "https://gitlab.com/" + *namespacedRepo + "/-/merge_requests?scope=all&state=merged&target_branch=" + to
	return &url, nil
}

func (gp *GitlabProvider) GetTask(taskId string) (*infrastructure.Task, error) {
	namespacedRepo, err := gp.namespacedRepo()
	if err != nil {
		return nil, err
	}

	gitlabClient, _ := gitlab.NewClient(gp.GitToken)
	project, _, err := gitlabClient.Projects.GetProject(*namespacedRepo, &gitlab.GetProjectOptions{})
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(strings.Replace(taskId, "!", "", 1))
	if err != nil {
		return nil, err
	}

	// Get the first page with projects.
	mergeRequest, _, err := gitlabClient.MergeRequests.GetMergeRequest(project.ID,
		id, &gitlab.GetMergeRequestsOptions{})
	if err != nil {
		return nil, err
	}

	var labelStrings []string
	labels := mergeRequest.Labels
	for _, label := range labels {
		labelStrings = append(labelStrings, label)
	}

	var fileChanges []string
	for _, change := range mergeRequest.Changes {
		fileChanges = append(fileChanges, change.NewPath)
	}

	gitTask := infrastructure.NewTask("!"+strconv.Itoa(mergeRequest.IID), "!"+strconv.Itoa(mergeRequest.IID), mergeRequest.Title, mergeRequest.WebURL, "@"+mergeRequest.Author.Username,
		mergeRequest.Author.WebURL, labelStrings, fileChanges)

	return gitTask, nil
}

func (gp *GitlabProvider) repoURL() (*string, error) {
	if gp.repo != nil {
		return gp.repo, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if Settings.Providers.Gitlab.GitRoot != "./" && Settings.Providers.Gitlab.GitRoot != "." && Settings.Providers.Gitlab.GitRoot != "" {
		cwd = filepath.Join(cwd, Settings.Providers.Gitlab.GitRoot)
	}

	r, err := git.PlainOpen(cwd)
	if err != nil {
		return nil, err
	}
	remotes, _ := r.Remotes()
	return &remotes[0].Config().URLs[0], nil
}

func (gp *GitlabProvider) namespacedRepo() (*string, error) {
	currentGitBAseUrl, err := gp.repoURL()
	if err != nil {
		return nil, err
	}

	*currentGitBAseUrl = strings.Replace(*currentGitBAseUrl, ".git", "", 1)
	namespacedRepoSliced := strings.Split(*currentGitBAseUrl, "gitlab.com/")
	if len(namespacedRepoSliced) <= 1 {
		namespacedRepoSliced = strings.Split(*currentGitBAseUrl, "gitlab.com:")
		if len(namespacedRepoSliced) <= 1 {
			Log.Fatalf("Unable to retrieve github repo/namespace from git origin")
		}
	}
	namespacedRepo := namespacedRepoSliced[1]

	return &namespacedRepo, nil
}
