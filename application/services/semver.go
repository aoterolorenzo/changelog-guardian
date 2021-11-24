package services

import (
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"regexp"
	"strconv"
)

type SemVerService struct {
}

func (svs *SemVerService) IsSemverValid() {

}

func (svs *SemVerService) CalculateNextVersion(categories []models.Category, versionToBump string) string {
	categoriesMap := make(map[models.Category]bool)

	for _, category := range categories {
		categoriesMap[category] = true
	}

	if categoriesMap[models.BREAKING_CHANGE] {
		return svs.BumpMajor(versionToBump)
	}

	if categoriesMap[models.ADDED] || categoriesMap[models.CHANGED] || categoriesMap[models.DEPRECATED] || categoriesMap[models.SECURITY] {
		return svs.BumpMinor(versionToBump)
	}

	return svs.BumpPatch(versionToBump)
}

func (svs *SemVerService) BumpPatch(versionToBump string) string {
	params := svs.getSemVerParams(versionToBump)

	newPatch, _ := strconv.Atoi(params["patch"])
	newPatch += 1
	params["patch"] = strconv.Itoa(newPatch)

	return params["major"] + "." + params["minor"] + "." + params["patch"]
}

func (svs *SemVerService) BumpMinor(versionToBump string) string {
	params := svs.getSemVerParams(versionToBump)

	newPatch, _ := strconv.Atoi(params["minor"])
	newPatch += 1
	params["minor"] = strconv.Itoa(newPatch)

	return params["major"] + "." + params["minor"] + ".0"
}

func (svs *SemVerService) BumpMajor(versionToBump string) string {
	params := svs.getSemVerParams(versionToBump)

	newPatch, _ := strconv.Atoi(params["major"])
	newPatch += 1
	params["major"] = strconv.Itoa(newPatch)

	return params["major"] + "." + "0.0"
}

func (svs *SemVerService) getSemVerParams(versionToBump string) map[string]string {
	var compRegEx = regexp.MustCompile(`^v?(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	match := compRegEx.FindStringSubmatch(versionToBump)

	paramsMap := make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}

	return paramsMap
}
