package cmd

import (
	"github.com/spf13/cobra"

	"github.com/motemen/sbx/lib/config"
	"github.com/motemen/sbx/lib/sbapi"
)

var projectCmd = &cobra.Command{
	Use:   "project <name>",
	Short: "Show project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		if optSession == "" {
			var err error
			optSession, err = config.GetSession(projectName)
			if err != nil {
				return err
			}
		}

		project, err := sbapi.GetProject(projectName, sbapi.WithSessionID(optSession))
		if err != nil {
			return err
		}

		return printResult(project, optJQQuery.Query)
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
