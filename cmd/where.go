package cmd

import (
	"fmt"

	"github.com/jclem/cobble/cobble/config"
	"github.com/spf13/cobra"
)

// whereCmd represents the run command
var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Print the location of the scaffolds directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := config.ScaffoldsDir()
		if err != nil {
			return err
		}

		_, err = fmt.Println(s)
		return err
	},
}

func init() {
	rootCmd.AddCommand(whereCmd)
}
