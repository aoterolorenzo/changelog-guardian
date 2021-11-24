package cli

import (
	"github.com/spf13/cobra"
	"gitlab.com/aoterocom/changelog-guardian/application/usecases"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Generates a new Release",
	Long: `Generates a new Release on the CHANGELOG
defining the type of bum following Semantic Versioning specification'`,
	Run: func(cmd *cobra.Command, args []string) {
		usecases.ReleaseCmd()
	},
}

func init() {
	regularCmd.AddCommand(releaseCmd)
}
