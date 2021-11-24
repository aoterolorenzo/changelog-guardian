package services

import (
	"fmt"
	"gitlab.com/aoterocom/changelog-guardian/application/helpers"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
)

type ChangelogMixer struct {
}

func NewChangelogMixer() *ChangelogMixer {
	return &ChangelogMixer{}
}

func (cm *ChangelogMixer) MergeChangelogs(changelog1 models.Changelog, changelog2 models.Changelog) models.Changelog {

	// Generate an unreleased: If there is one in local changelog (which should), just save the merge between both
	// If not, just save the changelog2 unreleased section
	var unreleased models.Release
	if cm.ChangelogContainsRelease(changelog1,
		*models.NewRelease("UNRELEASED", "", "", false, nil)) {
		unreleased = *cm.MergeReleases(changelog2.Releases[0], changelog1.Releases[0])
		changelog1.Releases = changelog1.Releases[1:]

	} else {
		unreleased = changelog2.Releases[0]
	}

	// Remove the unreleased section to start iterating to merge the releases
	changelog2.Releases = changelog2.Releases[1:]

	// Now, we just iterate on changelog2 releases, saving or merging each depending on if exist or not
	for i, release := range changelog2.Releases {
		if cm.ChangelogContainsRelease(changelog1, release) {
			changelog1.Releases[i] = *cm.MergeReleases(changelog1.Releases[i], changelog2.Releases[i])
		} else {
			changelog1.Releases = append(changelog1.Releases, release)
		}
	}

	// Finally, we append the unreleased section to the top changelog1
	// (array order is desc, so we must reverse it before adding)
	helpers.ReverseAny(changelog1.Releases)
	changelog1.Releases = append(changelog1.Releases, unreleased)
	helpers.ReverseAny(changelog1.Releases)

	return changelog1
}

func (cm *ChangelogMixer) MergeReleases(release1 models.Release, release2 models.Release) *models.Release {

	for category, section := range release1.Sections {
		for _, task := range section {
			// If task in both releases, we remove the one on release1, and we take the one in release2
			_, _, release2containsTask := cm.ReleaseContainsTask(release2, task)
			if !release2containsTask {
				release2.Sections[category] = append(release2.Sections[category], task)
			}

		}
	}

	return &release2
}

func (cm *ChangelogMixer) ReleaseContainsTask(release models.Release, task models.Task) (*models.Category, *models.Task, bool) {
	for category, tasksInSection := range release.Sections {
		for _, taskInRelease := range tasksInSection {
			if task.ID == taskInRelease.ID && category != models.REMOVED {
				return &category, &task, true
				// If the last entry of the task is removed, it's not contained since then
			} else if task.ID == taskInRelease.ID && category == models.REMOVED {
				return nil, nil, false
			}
		}

	}
	return nil, nil, false
}

func (cm *ChangelogMixer) ChangelogContainsRelease(changelog models.Changelog, rel models.Release) bool {
	for _, release := range changelog.Releases {
		if release.Version == rel.Version {
			return true
		}
	}
	return false
}

func (cm *ChangelogMixer) ChangelogContainsTask(changelog models.Changelog, task models.Task) (*models.Category, *models.Task, bool) {
	for _, release := range changelog.Releases {
		for category, sectionTasks := range release.Sections {
			for _, taskInChangelog := range sectionTasks {
				fmt.Println(taskInChangelog)
				fmt.Println(taskInChangelog.Category)
				if task.ID == taskInChangelog.ID && category != models.REMOVED {
					return &category, &task, true
					// If the last entry of the task is removed, it's not contained since then
				} else if task.ID == taskInChangelog.ID && category == models.REMOVED {
					return nil, nil, false
				}
			}
		}
	}
	return nil, nil, false
}
