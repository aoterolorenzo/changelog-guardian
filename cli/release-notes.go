package cli

import (
	"gitlab.com/aoterocom/changelog-guardian/application/usecases"

	"github.com/spf13/cobra"
)

// releaseNotesCmd represents the release-notes command
var releaseNotesCmd = &cobra.Command{
	Use:   "release-notes",
	Short: "Print version release notes",
	Long:  `Print the release notes of the last or given version.`,
	Run: func(cmd *cobra.Command, args []string) {
		usecases.ReleaseNotesCmd(cmd, args)
	},
}

func init() {
	RegularCmd.AddCommand(releaseNotesCmd)
	releaseNotesCmd.Flags().StringP("version", "v", "", "Version")
	releaseNotesCmd.Flags().StringP("output-file", "o", "", "Output file")
	releaseNotesCmd.Flags().BoolP("echo", "e", false, "Echo Release Notes on screen")
}
