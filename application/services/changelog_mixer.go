package services

import (
	"fmt"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/helpers"
	"sort"
	"strings"
	"time"
)

type ChangelogMixer struct {
}

func NewChangelogMixer() *ChangelogMixer {
	return &ChangelogMixer{}
}

func (cm *ChangelogMixer) MergeChangelogs(changelog1 models.Changelog, changelog2 models.Changelog) models.Changelog {

	// If one is empty, return the other
	if len(changelog1.Releases) == 0 {
		return changelog2
	} else if len(changelog2.Releases) == 0 {
		return changelog1
	}
	// Generate an unreleased: If there is one in local changelog (which should), just save the merge between both
	// If not, just save the changelog2 unreleased section
	var unreleased models.Release
	if cm.ChangelogContainsRelease(changelog1,
		*models.NewRelease("Unreleased", "", "", false, nil)) {

		if cm.ChangelogContainsRelease(changelog2,
			*models.NewRelease("Unreleased", "", "", false, nil)) {
			unreleased = *cm.MergeReleases(changelog2.Releases[0], changelog1.Releases[0])
			// Remove the unreleased section to start iterating to merge the releases
			changelog2.Releases = changelog2.Releases[1:]
		} else {
			unreleased = changelog1.Releases[0]
		}

		changelog1.Releases = changelog1.Releases[1:]

	} else {
		unreleased = changelog2.Releases[0]
	}

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

	orderedChangelog := cm.orderChangelogReleasesByDate(changelog1)
	return cm.parseRemovals(orderedChangelog)
}

func (cm *ChangelogMixer) MergeReleases(release1 models.Release, release2 models.Release) *models.Release {

	for category, section := range release1.Sections {
		helpers.ReverseAny(section)
		for _, task := range section {
			// If task in both releases, we remove the one on release1, and we take the one in release2
			_, _, release2containsTask := cm.ReleaseContainsTask(release2, task)
			if !release2containsTask {
				if release2.Sections == nil {
					release2.Sections = make(map[models.Category][]models.Task)
				}

				helpers.ReverseAny(release2.Sections[category])
				release2.Sections[category] = append(release2.Sections[category], task)
				helpers.ReverseAny(release2.Sections[category])
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

func (cm *ChangelogMixer) parseRemovals(changelog models.Changelog) models.Changelog {
	for _, release := range changelog.Releases {
		for _, task := range release.Sections[models.REMOVED] {

			// Research for a non-REMOVED task with same ID or title IN SAME RELEASE
			for catT, sectionT := range release.Sections {
				if catT == models.REMOVED {
					continue
				}

				for _, taskT := range sectionT {
					fmt.Println(task.Title)
					fmt.Println(taskT.Title)
					if task.ID == taskT.ID || task.Title == taskT.Title ||
						strings.ToLower(task.Title) == "revert \""+strings.ToLower(taskT.Title)+"\"" {
						// REMOVE both
						release.Sections[models.REMOVED] = removeTaskFromRelease(release.Sections[models.REMOVED], task)
						if len(release.Sections[models.REMOVED]) == 0 {
							delete(release.Sections, models.REMOVED)
						}

						release.Sections[catT] = removeTaskFromRelease(release.Sections[catT], taskT)
						if len(release.Sections[catT]) == 0 {
							delete(release.Sections, catT)
						}
					}
				}
			}
		}
	}

	return changelog
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

func (cm *ChangelogMixer) orderChangelogReleasesByDate(changelog models.Changelog) models.Changelog {

	sort.Slice(changelog.Releases, func(i, j int) bool {

		if changelog.Releases[j].Version == "Unreleased" {
			return false
		}

		layout := "02-01-2006"
		date1, _ := time.Parse(layout, changelog.Releases[i].Date)
		date2, _ := time.Parse(layout, changelog.Releases[j].Date)
		return date1.Unix() > date2.Unix()
	})

	return changelog
}

func removeTaskFromRelease(s []models.Task, r models.Task) []models.Task {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
