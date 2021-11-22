package providers

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	infrastructure "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"os"
	"time"
)

type GitProvider struct {
}

func NewGitController() *GitProvider {
	return &GitProvider{}
}

func (gc *GitProvider) GetReleases(repo *string) (*[]infrastructure.Release, error) {

	var path string
	var err error
	if repo == nil {
		path, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	} else {
		path = *repo
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

	return &tags, nil
}

func (gc *GitProvider) GetTasks(from *time.Time, to *time.Time, repo *string, targetBranch string) (*[]infrastructure.Task, error) {
	//TODO: Implement GetTasks method
	return nil, nil
}
