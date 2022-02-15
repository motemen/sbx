package cmd

import (
	"github.com/spf13/cobra"

	"github.com/motemen/sbx/lib/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config related commands",
}

var configShowCmd = &cobra.Command{
	Use:   "show <project>",
	Short: "Show configuration for project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		projectConf, err := config.GetProjectConfig(projectName)
		cobra.CheckErr(err)

		err = printResult(projectConf, optJQQuery.Query)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configShowCmd)
}
