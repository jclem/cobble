package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/jclem/cobble/cobble/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cobble",
	Short:   "Cobble cobbles together projects from task definitions",
	Version: Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Ensures that usage isn't printed for errors such as non-zero exits.
		// SEE: https://github.com/spf13/cobra/issues/340#issuecomment-378726225
		cmd.SilenceUsage = true
	},
}

func Execute() {
	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		cobra.CheckErr(config.InitConfig())
	})
	rootCmd.PersistentFlags().String("config", "", fmt.Sprintf("config file (default is %s)", config.DefaultConfigPath))
	cobra.CheckErr(viper.BindPFlag(config.ConfigFileFlag, rootCmd.PersistentFlags().Lookup("config")))
}
