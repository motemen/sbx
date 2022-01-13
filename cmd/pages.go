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
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		if optSession == "" {
			var err error
			optSession, err = config.GetSession(projectName)
			if err != nil {
				return err
			}
		}

		pages, err := sbapi.ListPages(
			projectName,
			sbapi.WithSessionID(optSession),
			sbapi.WithLimit(limit),
		)
		if err != nil {
			return err
		}

		return printResult(pages, optJQQuery.Query)
	},
}

var limit uint

func init() {
	rootCmd.AddCommand(pagesCmd)

	pagesCmd.Flags().UintVarP(&limit, "limit", "L", 100, "limit")
}
