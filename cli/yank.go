package cli

import (
	"gitlab.com/aoterocom/changelog-guardian/application/usecases"

	"github.com/spf13/cobra"
)

// yankCmd represents the yank command
var yankCmd = &cobra.Command{
	Use:   "yank",
	Short: "Yank release",
	Long: `Yank a release and move its tasks to the immediately after 
release (or unrelease it no more releases are present.`,
	Run: func(cmd *cobra.Command, args []string) {
		usecases.YankCmd(cmd, args)
	},
}

func init() {
	RegularCmd.AddCommand(yankCmd)
	yankCmd.Flags().StringP("version", "v", "", "Version to yank")
}
