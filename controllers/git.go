package controllers

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	models "gitlab.com/aoterocom/changelog-guardian/models/git"
	"os"
)

type GitController struct {
}

func NewGitController() *GitController {
	return &GitController{}
}

func (gc *GitController) getTags() (*[]models.Tag, error) {
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

	var tags []models.Tag

	err = gitTags.ForEach(func(t *plumbing.Reference) error {
		revHash, err := r.ResolveRevision(plumbing.Revision(t.Name()))
		if err != nil {
			return nil
		}

		commit, err := r.CommitObject(*revHash)
		if err != nil {
			return nil
		}

		tag := models.NewTag(t.Name().Short(), models.Hash(revHash.String()), commit.Author.When)
		tags = append(tags, *tag)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &tags, nil
}

func (gc *GitController) getLastTag() (*models.Tag, error) {
	tags, err := gc.getTags()
	if err != nil {
		return nil, err
	}

	return &(*tags)[len(*tags)-1], err
}
