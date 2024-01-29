package providers

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v41/github"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	application "gitlab.com/aoterocom/changelog-guardian/application/models"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/helpers"
	infrastructure "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"golang.org/x/oauth2"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type GithubPRProvider struct {
	GitToken string
	repo     *string
}

func init() {
	cwd, _ := os.Getwd()
	_ = godotenv.Load(cwd + "/vars.env")
}

func NewGithubPRProvider(repo *string) *GithubProvider {
	return &GithubProvider{
		GitToken: os.Getenv("GITHUB_TOKEN"),
		repo:     repo,
	}
}

func (gp *GithubPRProvider) GetReleases(from *time.Time, to *time.Time) (*[]infrastructure.Release, error) {

	namespacedRepo, err := gp.namespacedRepo()
	if err != nil {
		return nil, err
	}

	// Prepare client auth
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gp.GitToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)

	var timeQueryStr string
	layout := "2006-01-02T15:04:05"
	if from == nil && to == nil {
		timeQueryStr = ""
	} else if from == nil {
		timeQueryStr = " created:<" + to.Format(layout)
	} else if to == nil {
		timeQueryStr = " created:>" + from.Format(layout)
	} else {
		timeQueryStr = " created:" + from.Format(layout) + ".." + to.Format(layout)
	}

	query := "type:pr is:merged repo:" + *namespacedRepo + " base:" + Settings.Providers.GithubPRs.TargetBranch +
		timeQueryStr + " " + Settings.Providers.GithubPRs.GHReleaseSearch
	opts := &github.SearchOptions{
		Sort:      "",
		Order:     "",
		TextMatch: false,
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 0,
		},
	}

	var gitReleases []infrastructure.Release
	var pullRequests []*github.Issue

	for {
		// Avoid rate limits reaching :_(
		time.Sleep(1 * time.Second)
		// Get the first page with projects.
		mergeRequestsRetrieved, resp, err := githubClient.Search.Issues(ctx, query, opts)
		if err != nil {
			return nil, err
		}

		pullRequests = append(pullRequests, mergeRequestsRetrieved.Issues...)

		// Exit the loop when we've seen all pages.
		if resp.NextPage == 0 {
			break
		}

		// Update the page number to get the next page.
		opts.Page = resp.NextPage
		opts.ListOptions.Page = resp.NextPage
	}

	for _, pullRequest := range pullRequests {

		var compRegEx = regexp.MustCompile(Settings.Providers.GithubPRs.VersionRegex)
		match := compRegEx.FindStringSubmatch(*pullRequest.Title)
		if len(match) == 0 {
			return nil, errors.New("version not matching")
		}

		fmt.Println(match)

		pullRequestId := *pullRequest.ID
		listOpts := &github.ListOptions{}
		var prFiles []*github.CommitFile
		for {
			// Avoid rate limits reaching :_(
			time.Sleep(300 * time.Millisecond)
			files, resp, _ := githubClient.PullRequests.ListFiles(ctx, gp.getOrg(*namespacedRepo), gp.getRepo(*namespacedRepo), int(pullRequestId), listOpts)
			prFiles = append(prFiles, files...)
			// Exit the loop when we've seen all pages.
			if resp.NextPage == 0 {
				break
			}

			// Update the page number to get the next page.
			listOpts.Page = resp.NextPage
		}

		var labelStrings []string
		labels := pullRequest.Labels
		for _, label := range labels {
			labelStrings = append(labelStrings, *label.Name)
		}

		var fileChanges []string
		for _, file := range prFiles {
			fileChanges = append(fileChanges, file.GetFilename())
		}

		releaseLink := *pullRequest.HTMLURL
		gitReleases = append(gitReleases, *infrastructure.NewRelease(*pullRequest.Title, infrastructure.Hash(rune(*pullRequest.Number)), *pullRequest.CreatedAt, releaseLink))
	}

	helpers.ReverseAny(gitReleases)
	return &gitReleases, nil
}

func (gp *GithubPRProvider) GetTasks(from *time.Time, to *time.Time, targetBranch string) (*[]infrastructure.Task, error) {
	return NewGithubProvider(gp.repo).GetTasks(from, to, targetBranch)
}

func (gp *GithubPRProvider) DefineCategory(task infrastructure.Task) application.Category {
	var category = application.ADDED

	if strings.HasPrefix(strings.ToLower(task.Title), "revert") {
		return application.REMOVED
	}

	for _, label := range task.Labels {
		switch label {
		case Settings.Providers.Github.Labels[application.ADDED]:
			category = application.ADDED
			break
		case Settings.Providers.Github.Labels[application.FIXED]:
			category = application.FIXED
			break
		case Settings.Providers.Github.Labels[application.REFACTOR]:
			category = application.REFACTOR
			break
		case Settings.Providers.Github.Labels[application.DEPRECATED]:
			category = application.DEPRECATED
			break
		case Settings.Providers.Github.Labels[application.CHANGED]:
			category = application.CHANGED
			break
		case Settings.Providers.Github.Labels[application.DEPENDENCIES]:
			category = application.DEPENDENCIES
			break
		case Settings.Providers.Github.Labels[application.DOCUMENTATION]:
			category = application.DOCUMENTATION
			break
		case Settings.Providers.Github.Labels[application.REMOVED]:
			category = application.REMOVED
			break
		case Settings.Providers.Github.Labels[application.SECURITY]:
			category = application.SECURITY
			break
		case Settings.Providers.Github.Labels[application.BREAKING_CHANGE]:
			category = application.BREAKING_CHANGE
			return category
		default:
			break
		}
	}

	return category
}

func (gp *GithubPRProvider) repoURL() (*string, error) {
	if gp.repo != nil {
		return gp.repo, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if Settings.Providers.GithubPRs.GitRoot != "./" && Settings.Providers.GithubPRs.GitRoot != "." && Settings.Providers.GithubPRs.GitRoot != "" {
		cwd = filepath.Join(cwd, Settings.Providers.GithubPRs.GitRoot)
	}

	r, err := git.PlainOpen(cwd)
	if err != nil {
		return nil, err
	}
	remotes, _ := r.Remotes()
	return &remotes[0].Config().URLs[0], nil
}

func (gp *GithubPRProvider) namespacedRepo() (*string, error) {
	currentGitBAseUrl, err := gp.repoURL()
	if err != nil {
		return nil, err
	}

	*currentGitBAseUrl = strings.Replace(*currentGitBAseUrl, ".git", "", 1)
	namespacedRepoSliced := strings.Split(*currentGitBAseUrl, "github.com/")
	if len(namespacedRepoSliced) <= 1 {
		namespacedRepoSliced = strings.Split(*currentGitBAseUrl, "github.com:")
		if len(namespacedRepoSliced) <= 1 {
			Log.Fatalf("Unable to retrieve github repo/namespace from git origin")
		}
	}
	namespacedRepo := namespacedRepoSliced[1]

	return &namespacedRepo, nil
}

func (gp *GithubPRProvider) getRepo(namespacedRepo string) string {
	slice := strings.Split(namespacedRepo, "/")
	return slice[1]
}

func (gp *GithubPRProvider) getOrg(namespacedRepo string) string {
	slice := strings.Split(namespacedRepo, "/")
	return slice[0]
}
