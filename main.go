package main

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/aoterocom/changelog-guardian/cli"
)

var (
	//go:embed VERSION
	version string
)

func main() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints version info",
		Long:  `Prints Changelog Guardian version and build info`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", version)
		},
	}
	cli.RegularCmd.AddCommand(versionCmd)
	cli.Execute()
}
