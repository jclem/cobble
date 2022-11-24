package cmd

import (
	"strings"

	"github.com/jclem/cobble/cobble/config"
	"github.com/jclem/cobble/cobble/runner"
	"github.com/jclem/cobble/cobble/task"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:     "run (SCAFFOLD...)",
	Aliases: []string{"r"},
	Short:   "Run one or more scaffolds",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workingDir, err := config.WorkingDir()
		if err != nil {
			return err
		}

		scaffoldsDir, err := config.ScaffoldsDir()
		if err != nil {
			return err
		}

		r, err := runner.NewWithOpts(
			runner.WithScaffoldsDir(scaffoldsDir),
			runner.WithWorkingDir(workingDir),
		)
		if err != nil {
			return err
		}

		return r.Run(args...)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		scaffolds, err := task.List()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var matchedScaffolds []string

		for _, scaffold := range scaffolds {
			if strings.HasPrefix(scaffold, toComplete) {
				matchedScaffolds = append(matchedScaffolds, scaffold)
			}
		}

		return matchedScaffolds, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	scaffoldsDirArg := "scaffolds-dir"
	workingDirArg := "working-dir"
	runCmd.Flags().StringP(scaffoldsDirArg, "d", "", "Directory containing scaffolds")
	runCmd.Flags().StringP(workingDirArg, "w", ".", "Working directory in which to execute scaffolds")

	cobra.CheckErr(viper.BindPFlag(config.ScaffoldsDirFlag, runCmd.Flags().Lookup(scaffoldsDirArg)))
	cobra.CheckErr(viper.BindPFlag(config.WorkingDirFlag, runCmd.Flags().Lookup(workingDirArg)))
}
