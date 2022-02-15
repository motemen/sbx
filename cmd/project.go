package cmd

import (
	"github.com/spf13/cobra"

	"github.com/motemen/sbx/lib/sbapi"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project related commands",
}

var projectShowCmd = &cobra.Command{
	Use:   "show <project>",
	Short: "Show project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		opts, err := buildApiOptions(projectName)
		cobra.CheckErr(err)

		project, err := sbapi.GetProject(projectName, opts...)
		cobra.CheckErr(err)

		err = printResult(project, optJQQuery.Query)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)

	projectCmd.AddCommand(projectShowCmd)
}
