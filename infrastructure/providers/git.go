package providers

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	application "gitlab.com/aoterocom/changelog-guardian/application/models"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	infrastructure "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type GitProvider struct {
}

func NewGitProvider() *GitProvider {
	return &GitProvider{}
}

func (gc *GitProvider) GetReleases(from *time.Time, to *time.Time) (*[]infrastructure.Release, error) {

	var path string
	var err error

	path, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	if Settings.Providers.Git.GitRoot != "./" && Settings.Providers.Git.GitRoot != "." && Settings.Providers.Git.GitRoot != "" {
		path = filepath.Join(path, Settings.Providers.Git.GitRoot)
	}

	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	gitTags, err := r.Tags()
	if err != nil {
		return nil, err
	}

	var tags []infrastructure.Release

	err = gitTags.ForEach(func(t *plumbing.Reference) error {
		revHash, err := r.ResolveRevision(plumbing.Revision(t.Name()))
		if err != nil {
			return nil
		}

		commit, err := r.CommitObject(*revHash)
		if err != nil {
			return nil
		}

		tag := infrastructure.NewRelease(t.Name().Short(), infrastructure.Hash(revHash.String()), commit.Author.When, "")
		tags = append(tags, *tag)

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.SliceStable(tags, func(i, j int) bool {
		return tags[i].Time.Unix() < tags[j].Time.Unix()
	})

	return &tags, nil
}

func (gc *GitProvider) GetTasks(from *time.Time, to *time.Time, targetBranch string) (*[]infrastructure.Task, error) {
	var path string
	var err error

	path, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	if Settings.Providers.Git.GitRoot != "./" && Settings.Providers.Git.GitRoot != "." && Settings.Providers.Git.GitRoot != "" {
		path = filepath.Join(path, Settings.Providers.Git.GitRoot)
	}

	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(targetBranch),
	})

	opts := &git.LogOptions{
		Since: from,
		Until: to,
	}

	log, err := r.Log(opts)
	if err != nil {
		return nil, err
	}

	var gitTasks []infrastructure.Task
	var resultFiles []string
	err = log.ForEach(func(commit *object.Commit) error {
		files, err := commit.Files()
		if err != nil {
			return err
		}

		err = files.ForEach(func(file *object.File) error {
			resultFiles = append(resultFiles, file.Name)
			return nil
		})
		if err != nil {
			return err
		}

		gitTask := infrastructure.NewTask(commit.Hash.String()[:6], commit.Hash.String()[:6], strings.Split(commit.Message, "\n")[0], commit.Hash.String(), commit.Author.Name,
			"mailto:"+commit.Author.Email, nil, resultFiles)
		gitTask.Category = gc.DefineCategory(*gitTask)
		gitTasks = append(gitTasks, *gitTask)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &gitTasks, nil
}

func (gc *GitProvider) DefineCategory(task infrastructure.Task) application.Category {
	if strings.HasPrefix(strings.ToLower(task.Title), "revert") {
		return application.REMOVED
	}
	return application.ADDED
}

func (gc *GitProvider) GetTask(taskId string) (*infrastructure.Task, error) {
	var path string
	var err error

	path, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	if Settings.Providers.Git.GitRoot != "./" && Settings.Providers.Git.GitRoot != "." && Settings.Providers.Git.GitRoot != "" {
		path = filepath.Join(path, Settings.Providers.Git.GitRoot)
	}

	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(Settings.DevelopBranch),
	})

	commit, err := r.CommitObject(plumbing.NewHash(taskId))
	if err != nil {
		return nil, err
	}

	files, err := commit.Files()
	if err != nil {
		return nil, err
	}

	var resultFiles []string
	err = files.ForEach(func(file *object.File) error {
		resultFiles = append(resultFiles, file.Name)
		return nil
	})
	if err != nil {
		return nil, err
	}

	gitTask := infrastructure.NewTask(commit.Hash.String()[:6], commit.Hash.String()[:6], strings.Split(commit.Message, "\n")[0], commit.Hash.String(), commit.Author.Name,
		"mailto:"+commit.Author.Email, nil, resultFiles)

	return gitTask, nil
}

func (gc *GitProvider) ReleaseURL(from *string, to string) (*string, error) {
	url := ""
	return &url, nil
}
