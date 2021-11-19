package providers

import (
	"github.com/go-git/go-git/v5"
	"github.com/joho/godotenv"
	"github.com/xanzy/go-gitlab"
	models2 "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"os"
	"strings"
	"time"
)

type GitlabProvider struct {
	gitToken string
}

func init() {
	cwd, _ := os.Getwd()
	_ = godotenv.Load(cwd + "/vars.env")
}

func NewGitlabProvider() *GitlabProvider {
	return &GitlabProvider{
		gitToken: os.Getenv("GITLAB_TOKEN"),
	}
}

func (glc *GitlabProvider) GetReleases() ([]*models2.Release, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	r, err := git.PlainOpen(cwd)
	if err != nil {
		return nil, err
	}
	remotes, _ := r.Remotes()
	currentGitBAseUrl := remotes[0].Config().URLs[0]
	gitlabClient, _ := gitlab.NewClient(glc.gitToken, gitlab.WithBaseURL(currentGitBAseUrl))

	gitlabProject := strings.Replace(currentGitBAseUrl, "https://gitlab.com/", "", 1)
	gitlabProject = strings.Replace(gitlabProject, ".git", "", 1)

	releases, _, _ := gitlabClient.Releases.ListReleases(gitlabProject, nil)
	if err != nil {
		return nil, err
	}

	var gitReleases []*models2.Release
	for _, release := range releases {
		gitReleases = append(gitReleases, models2.NewRelease(release.Name, models2.Hash(release.Commit.ID), *release.ReleasedAt))
	}
	return gitReleases, nil
}

func (glc *GitlabProvider) GetTasks(from *time.Time, to *time.Time, targetBranch string) (*[]models2.Task, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	r, err := git.PlainOpen(cwd)
	if err != nil {
		return nil, err
	}
	remotes, _ := r.Remotes()
	currentGitBAseUrl := remotes[0].Config().URLs[0]

	listMergeRequestOptions := &gitlab.ListMergeRequestsOptions{CreatedAfter: from, CreatedBefore: to}

	gitlabClient, _ := gitlab.NewClient(glc.gitToken, gitlab.WithBaseURL(currentGitBAseUrl))

	gitlabProject := strings.Replace(currentGitBAseUrl, "https://gitlab.com/", "", 1)
	gitlabProject = strings.Replace(gitlabProject, ".git", "", 1)

	var gitTasks []models2.Task
	mergeRequests, _, err := gitlabClient.MergeRequests.ListMergeRequests(listMergeRequestOptions)
	if err != nil {
		return nil, err
	}

	for _, mergeRequest := range mergeRequests {

		if mergeRequest.TargetBranch == targetBranch && mergeRequest.MergeStatus == "Merged" {
			var labelStrings []string
			labels := mergeRequest.Labels
			for _, label := range labels {
				labelStrings = append(labelStrings, label)
			}

			var fileChanges []string
			for _, change := range mergeRequest.Changes {
				fileChanges = append(fileChanges, change.NewPath)
			}

			gitTask := models2.NewTask(mergeRequest.Title, mergeRequest.WebURL, mergeRequest.Author.Username,
				mergeRequest.Author.WebURL, labelStrings, fileChanges)
			gitTasks = append(gitTasks, *gitTask)
		}
	}

	return &gitTasks, nil
}
