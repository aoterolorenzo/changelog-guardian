package cli

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	. "gitlab.com/aoterocom/changelog-guardian/config"
	"strconv"
)

func PreCommandChecks(cmd *cobra.Command, args []string) {
	argChangelogPath := cmd.Flag("changelog-path").Value.String()
	if argChangelogPath != "" {
		Settings.ChangelogPath = argChangelogPath
	}
	argConfigFile := cmd.Flag("config").Value.String()
	if argConfigFile != "" {
		Settings.CGConfigPath = argConfigFile
	}

	argSilent, _ := strconv.ParseBool(cmd.Flag("silent").Value.String())
	if argSilent {
		Log.SetLevel(log.ErrorLevel)
	}
}
