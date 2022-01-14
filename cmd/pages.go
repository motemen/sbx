package cmd

import (
	"github.com/spf13/cobra"

	"github.com/motemen/sbx/lib/config"
	"github.com/motemen/sbx/lib/sbapi"
)

type PagesResponse struct {
	ProjectName string        `json:"projectName"`
	Skip        int           `json:"skip"`
	Limit       int           `json:"limit"`
	Count       int           `json:"count"`
	Pages       []interface{} `json:"pages"`
}

type ErrorResponse struct {
	Name    string      `json:"name"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

var pagesCmd = &cobra.Command{
	Use:   "pages <project>",
	Short: "List pages in project",
	Args:  cobra.ExactArgs(1),
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
	rootCmd.AddCommand(pagesCmd)

	pagesCmd.Flags().UintVarP(&limit, "limit", "L", 100, "limit")
}
