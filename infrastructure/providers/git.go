package providers

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	models2 "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"os"
	"time"
)

type GitProvider struct {
}

func NewGitController() *GitProvider {
	return &GitProvider{}
}

func (gc *GitProvider) GetReleases() (*[]models2.Release, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	r, err := git.PlainOpen(cwd)
	if err != nil {
		return nil, err
	}

	gitTags, err := r.Tags()
	if err != nil {
		return nil, err
	}

	var tags []models2.Release

	err = gitTags.ForEach(func(t *plumbing.Reference) error {
		revHash, err := r.ResolveRevision(plumbing.Revision(t.Name()))
		if err != nil {
			return nil
		}

		commit, err := r.CommitObject(*revHash)
		if err != nil {
			return nil
		}

		tag := models2.NewRelease(t.Name().Short(), models2.Hash(revHash.String()), commit.Author.When)
		tags = append(tags, *tag)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &tags, nil
}

func (gc *GitProvider) GetTasks(from *time.Time, to *time.Time, targetBranch string) (*[]models2.Task, error) {
	//TODO: Implement GetTasks method
	return nil, nil
}
