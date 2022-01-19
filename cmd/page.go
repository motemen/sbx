package cmd

import (
	"github.com/spf13/cobra"

	"github.com/motemen/sbx/lib/config"
	"github.com/motemen/sbx/lib/sbapi"
)

var pageCmd = &cobra.Command{
	Use:   "page",
	Short: "Page related commands",
}

var pageListCmd = &cobra.Command{
	Use:  "list [--limit <limit>] <project>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		if optSession == "" {
			var err error
			optSession, err = config.GetSession(projectName)
			cobra.CheckErr(err)
		}

		pages, err := sbapi.ListPages(
			projectName,
			sbapi.WithSessionID(optSession),
			sbapi.WithLimit(limit),
		)
		cobra.CheckErr(err)

		err = printResult(pages, optJQQuery.Query)
		cobra.CheckErr(err)
	},
}

var limit uint

func init() {
	rootCmd.AddCommand(pageCmd)

	pageCmd.AddCommand(pageListCmd)
	pageListCmd.Flags().UintVarP(&limit, "limit", "L", 100, "limit")
}
