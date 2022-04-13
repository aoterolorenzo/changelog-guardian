package cli

import (
	"github.com/spf13/cobra"
	"gitlab.com/aoterocom/changelog-guardian/application/usecases"
)

// releaseCmd represents the release command
var calculateReleaseCmd = &cobra.Command{
	Use:   "calculate-release",
	Short: "Calculates next release version",
	Long: `Calculates next release version
defining the type of bum following Semantic Versioning specification'`,
	Run: func(cmd *cobra.Command, args []string) {
		PreCommandChecks(cmd, args)
		usecases.CalculateReleaseCmd(cmd, args)
	},
}

func init() {
	RegularCmd.AddCommand(calculateReleaseCmd)

	calculateReleaseCmd.Flags().BoolP("patch", "p", false, "Patch Release")
	calculateReleaseCmd.Flags().BoolP("minor", "m", false, "Minor Release")
	calculateReleaseCmd.Flags().BoolP("major", "M", false, "Major Release")
	calculateReleaseCmd.Flags().String("pre", "", "Pre-release string (semver)")
	calculateReleaseCmd.Flags().String("build", "", "Build metadata (semver)")
	calculateReleaseCmd.Flags().BoolP("skip-update", "", false, "Skip CHANGELOG.md update before executing the command")

}
