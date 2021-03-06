package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/motemen/sbx/lib/sbapi"
)

var apiCmd = &cobra.Command{
	Use:  "api <path>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		if !strings.HasPrefix(path, "api/") {
			path = "api/" + path
		}

		// XXX: Guess project name
		// path should be shape of api/<category>/<project>
		// https://scrapbox.io/scrapboxlab/Scrapbox_API%E3%81%AE%E4%B8%80%E8%A6%A7
		var projectName string
		if parts := strings.Split(path, "/"); len(parts) > 2 {
			projectName = parts[2]
		}

		opts, err := buildApiOptions(projectName)
		cobra.CheckErr(err)

		var v interface{}
		err = sbapi.RequestJSON("/"+path, &v, opts...)
		cobra.CheckErr(err)

		if b, ok := v.([]byte); ok {
			fmt.Print(string(b))
			return
		}

		err = printResult(v, optJQQuery.Query)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
