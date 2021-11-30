package providers

import (
	"github.com/go-git/go-git/v5"
	"github.com/joho/godotenv"
	"github.com/xanzy/go-gitlab"
	application "gitlab.com/aoterocom/changelog-guardian/application/models"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	infrastructure "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"os"
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
	currentGitBAseUrl, err := gp.repoURL()
	if err != nil {
		return nil, err
	}

	gitlabProjectWebUrl := strings.Replace(*currentGitBAseUrl, ".git", "", 1)
	namespacedRepoSlice := strings.Split(gitlabProjectWebUrl, "gitlab.com/")
	var namespacedRepo string
	if len(namespacedRepoSlice) > 1 {
		namespacedRepo = namespacedRepoSlice[1]
	} else {
		Log.Fatalf("Error retrieving Gitlab repository from git origin")
	}

	gitlabClient, _ := gitlab.NewClient(gp.GitToken)
	project, _, err := gitlabClient.Projects.GetProject(namespacedRepo, &gitlab.GetProjectOptions{})
	if err != nil {
		return nil, err
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
		releaseLink := gitlabProjectWebUrl + "/-/releases/" + release.Name
		gitReleases = append(gitReleases, *infrastructure.NewRelease(release.Name, infrastructure.Hash(release.Commit.ID), *release.ReleasedAt, releaseLink))
	}
	return &gitReleases, nil
}

func (gp *GitlabProvider) GetTasks(from *time.Time, to *time.Time, targetBranch string) (*[]infrastructure.Task, error) {

	currentGitBAseUrl, err := gp.repoURL()
	if err != nil {
		return nil, err
	}

	gitlabProjectName := strings.Replace(*currentGitBAseUrl, "https://gitlab.com/", "", 1)
	gitlabProjectName = strings.Replace(gitlabProjectName, ".git", "", 1)

	gitlabClient, _ := gitlab.NewClient(gp.GitToken)
	project, _, err := gitlabClient.Projects.GetProject(gitlabProjectName, &gitlab.GetProjectOptions{})
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
	repo, err := gp.repoURL()
	if err != nil {
		return nil, err
	}
	gitlabProjectName := strings.Replace(*repo, "https://gitlab.com/", "", 1)
	gitlabProjectName = strings.Replace(gitlabProjectName, ".git", "", 1)
	if from != nil {
		url := "https://gitlab.com/" + gitlabProjectName + "/-/compare/" + *from + "..." + to
		return &url, nil
	}

	url := "https://gitlab.com/" + gitlabProjectName + "/-/merge_requests?scope=all&state=merged&target_branch=" + to
	return &url, nil
}

func (gp *GitlabProvider) GetTask(taskId string) (*infrastructure.Task, error) {
	currentGitBAseUrl, err := gp.repoURL()
	if err != nil {
		return nil, err
	}

	gitlabProjectWebUrl := strings.Replace(*currentGitBAseUrl, ".git", "", 1)
	namespacedRepoSlice := strings.Split(gitlabProjectWebUrl, "gitlab.com/")
	var namespacedRepo string
	if len(namespacedRepoSlice) > 1 {
		namespacedRepo = namespacedRepoSlice[1]
	} else {
		Log.Fatalf("Error retrieving Gitlab repository from git origin")
	}

	gitlabClient, _ := gitlab.NewClient(gp.GitToken)
	project, _, err := gitlabClient.Projects.GetProject(namespacedRepo, &gitlab.GetProjectOptions{})
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
	r, err := git.PlainOpen(cwd)
	if err != nil {
		return nil, err
	}
	remotes, _ := r.Remotes()
	return &remotes[0].Config().URLs[0], nil
}
