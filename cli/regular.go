package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/aoterocom/changelog-guardian/application/usecases"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// regularCmd represents the base command when called without any subcommands
var regularCmd = &cobra.Command{
	Use:   "Changelog Guardian",
	Short: "Keep you're changelog safe",
	Long:  `Keep you're changelog safe and punish those who dare to manually edit it`,

	Run: func(cmd *cobra.Command, args []string) {
		usecases.RegularCmd(cmd, args)
	},
}

func Execute() {
	cobra.CheckErr(regularCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	regularCmd.PersistentFlags().String("template", "", "CHANGELOG template")
	regularCmd.PersistentFlags().String("output-template", "", "Output CHANGELOG template")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".CLogger" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".clg.yml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
