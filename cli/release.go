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
		PreCommandChecks(cmd, args)
		usecases.ReleaseCmd(cmd, args)
	},
}

func init() {
	RegularCmd.AddCommand(releaseCmd)

	releaseCmd.Flags().BoolP("patch", "p", false, "Patch Release")
	releaseCmd.Flags().BoolP("minor", "m", false, "Minor Release")
	releaseCmd.Flags().BoolP("major", "M", false, "Major Release")
	releaseCmd.Flags().BoolP("force", "f", false, "Force versioning")
	releaseCmd.Flags().StringP("version", "v", "", "Force versioning")
	releaseCmd.Flags().String("pre", "", "Pre-release string (semver)")
	releaseCmd.Flags().String("build", "", "Build metadata (semver)")
	releaseCmd.Flags().BoolP("skip-update", "", false, "Skip CHANGELOG.md update before executing the command")

}
