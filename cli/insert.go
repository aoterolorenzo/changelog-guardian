package cli

import (
	"github.com/spf13/cobra"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/application/usecases"
)

// releaseCmd represents the release command
var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Inserts a task in CHANGELOG",
	Long: `Inserts a new task in the UNRELEASED section of
the CHANGELOG'`,
	Run: func(cmd *cobra.Command, args []string) {
		usecases.InsertCmd(cmd, args)
	},
}

func init() {
	RegularCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringP("title", "t", "", "Task title")
	insertCmd.Flags().StringP("id", "i", "", "Task ID")
	insertCmd.Flags().StringP("link", "l", "", "Task link")
	insertCmd.Flags().StringP("author", "a", "", "Task author")
	insertCmd.Flags().String("authorLink", "", "Task author link")
	insertCmd.Flags().StringP("category", "c", string(models.ADDED), "Task category")
	insertCmd.Flags().BoolP("skip-autocompletion", "s", false, "Skip autocompletion from provider"+
		"Used to check the task data from it through the provided --id")

}
