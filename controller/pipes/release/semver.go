package middleware

import (
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"regexp"
)

type SemverReleasePipe struct {
}

func NewSemverReleasePipe() *SemverReleasePipe {
	return &SemverReleasePipe{}
}

func (nlm *SemverReleasePipe) Pipe(release *infra.Release) (*infra.Release, bool, error) {
	match, err := regexp.Match(`^v?(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`, []byte(release.Name))
	if err != nil {
		return nil, true, err
	}

	if match {
		return release, false, nil
	} else {
		return nil, true, nil
	}
}
