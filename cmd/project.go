package cmd

import (
	"github.com/spf13/cobra"

	"github.com/motemen/sbx/lib/config"
	"github.com/motemen/sbx/lib/sbapi"
)

var projectCmd = &cobra.Command{
	Use: "project",
}

var projectShowCmd = &cobra.Command{
	Use:   "show <project>",
	Short: "Show project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		if optSession == "" {
			var err error
			optSession, err = config.GetSession(projectName)
			cobra.CheckErr(err)
		}

		project, err := sbapi.GetProject(projectName, sbapi.WithSessionID(optSession))
		cobra.CheckErr(err)

		err = printResult(project, optJQQuery.Query)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)

	projectCmd.AddCommand(projectShowCmd)
}
