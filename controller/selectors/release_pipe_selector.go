package selectors

import (
	"github.com/pkg/errors"
	"gitlab.com/aoterocom/changelog-guardian/controller/interfaces"
	releasePipes "gitlab.com/aoterocom/changelog-guardian/controller/pipes/release"
)

func ReleasePipeSelector(providerStr string) (*interfaces.ReleasePipe, error) {
	switch providerStr {
	case "semver":
		prov := interfaces.ReleasePipe(releasePipes.NewSemverReleasePipe())
		return &prov, nil
	default:
		return nil, errors.Errorf("unknown release pipe " + providerStr)
	}
}
