package selectors

import (
	"github.com/pkg/errors"
	"gitlab.com/aoterocom/changelog-guardian/application/interfaces"
	"gitlab.com/aoterocom/changelog-guardian/application/themes"
)

func ChangelogTemplateSelector(themeStr string) (*interfaces.ChangelogService, error) {
	switch themeStr {
	case "markdown":
		srv := interfaces.ChangelogService(themes.NewMarkDownChangelogService())
		return &srv, nil
	case "stylish_markdown":
		srv := interfaces.ChangelogService(themes.NewStylishMarkDownChangelogService())
		return &srv, nil
	default:
		return nil, errors.Errorf("unknown provider " + themeStr)
	}
}
