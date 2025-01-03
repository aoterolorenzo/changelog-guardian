package pipes

import (
	"fmt"
	settings "gitlab.com/aoterocom/changelog-guardian/config"
	"gitlab.com/aoterocom/changelog-guardian/helpers"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"strings"
)

type InclusionsExclusionsTasksPipe struct {
}

func NewInclusionsExclusionsTasksPipe() *InclusionsExclusionsTasksPipe {
	return &InclusionsExclusionsTasksPipe{}
}

func (tf *InclusionsExclusionsTasksPipe) Filter(task *infra.Task) (*infra.Task, bool, error) {
	var addressedLabelInclusions = false
	var addressedLabelExclusions = true

	if helpers.SliceContainsString(settings.Settings.TasksPipesCfg.InclusionsExclusions.Labels.Inclusions, "*all") {
		addressedLabelInclusions = true
	}

	for _, label := range task.Labels {

		if !helpers.SliceContainsString(settings.Settings.TasksPipesCfg.InclusionsExclusions.Labels.Inclusions, "*all") &&
			helpers.SliceContainsString(settings.Settings.TasksPipesCfg.InclusionsExclusions.Labels.Inclusions, label) {
			addressedLabelInclusions = true
		}

		if helpers.SliceContainsString(settings.Settings.TasksPipesCfg.InclusionsExclusions.Labels.Exclusions, label) {
			addressedLabelExclusions = false
		}
	}

	var addressedPathInclusions = false
	var addressedPathExclusions = true

	if helpers.SliceContainsString(settings.Settings.TasksPipesCfg.InclusionsExclusions.Paths.Inclusions, "*all") {
		addressedPathInclusions = true
	}

	for _, filePath := range task.Files {

		if !helpers.SliceContainsString(settings.Settings.TasksPipesCfg.InclusionsExclusions.Paths.Inclusions, "*all") {
			for _, path := range settings.Settings.TasksPipesCfg.InclusionsExclusions.Paths.Inclusions {
				path = strings.TrimSuffix(path, "/")
				split := strings.Split(filePath, "/")
				pathWithoutFile := strings.Join(split[:len(split)-1], "/")
				if path == pathWithoutFile {
					addressedPathInclusions = true
				}

			}
		}

		for _, path := range settings.Settings.TasksPipesCfg.InclusionsExclusions.Paths.Exclusions {
			path = strings.TrimSuffix(path, "/")
			split := strings.Split(filePath, "/")
			pathWithoutFile := strings.Join(split[:len(split)-1], "/")
			if path == pathWithoutFile {
				addressedPathExclusions = false
			}

		}

	}

	var addressedRegexpsInclusions = false
	var addressedRegexpsExclusions = true

	if helpers.SliceContainsString(settings.Settings.TasksPipesCfg.InclusionsExclusions.Regexps.Inclusions, "*all") {
		addressedRegexpsInclusions = true
	}

	if !helpers.SliceContainsString(settings.Settings.TasksPipesCfg.InclusionsExclusions.Regexps.Inclusions, "*all") &&
		helpers.StringMatchesRegexSlice(settings.Settings.TasksPipesCfg.InclusionsExclusions.Regexps.Inclusions, task.Name) {
		addressedRegexpsInclusions = true
	}

	if helpers.StringMatchesRegexSlice(settings.Settings.TasksPipesCfg.InclusionsExclusions.Labels.Exclusions, task.Name) {
		addressedRegexpsExclusions = false
	}

	if !addressedLabelInclusions {
		settings.Log.Debug("Task skipped: task label not between inclusion labels")
	} else if !addressedLabelExclusions {
		settings.Log.Debug("Task skipped: task label in exclusion labels")
	}

	if !addressedRegexpsInclusions {
		settings.Log.Debug(fmt.Sprintf("Task skipped: task name %s not between inclusion regexps", task.Title))
	} else if !addressedRegexpsExclusions {
		settings.Log.Debug(fmt.Sprintf("Task skipped: task name %s in exclusion regexps", task.Title))
	}

	if !addressedPathInclusions {
		settings.Log.Debug("Task skipped: task file path not between inclusion paths")
	} else if !addressedPathExclusions {
		settings.Log.Debug("Task skipped: task file path in exclusion paths")
	}

	if addressedLabelInclusions && addressedLabelExclusions &&
		addressedPathInclusions && addressedPathExclusions &&
		addressedRegexpsInclusions && addressedRegexpsExclusions {
		return task, false, nil
	}
	return nil, true, nil
}
