package selectors

import (
	"github.com/pkg/errors"
	"gitlab.com/aoterocom/changelog-guardian/infrastructure/interfaces"
	"gitlab.com/aoterocom/changelog-guardian/infrastructure/providers"
)

func ProviderSelector(providerStr string) (*interfaces.Provider, error) {
	switch providerStr {
	case "git":
		prov := interfaces.Provider(providers.NewGitProvider())
		return &prov, nil
	case "gitlab":
		prov := interfaces.Provider(providers.NewGitlabProvider(nil))
		return &prov, nil
	case "github":
		prov := interfaces.Provider(providers.NewGithubProvider(nil))
		return &prov, nil
	case "githubPrs":
		prov := interfaces.Provider(providers.NewGithubPRProvider(nil))
		return &prov, nil
	default:
		return nil, errors.Errorf("unknown provider " + providerStr)
	}
}
