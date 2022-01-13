package cmd

import (
	"net/http"
	"os"

	"github.com/itchyny/gojq"
	"github.com/spf13/cobra"

	"github.com/motemen/go-loghttp"
)

var rootCmd = &cobra.Command{
	Use:   "sbx",
	Short: "An unofficial client for Scrapbox",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	optJQQuery jqQueryFlag
	optSession string
)

type jqQueryFlag struct {
	*gojq.Query
}

func (q jqQueryFlag) String() string {
	if q.Query == nil {
		return ""
	}
	return q.Query.String()
}

func (q *jqQueryFlag) Set(s string) (err error) {
	q.Query, err = gojq.Parse(s)
	return err
}

func (q jqQueryFlag) Type() string {
	return "query"
}

func init() {
	rootCmd.PersistentFlags().VarP(&optJQQuery, "jq", "q", "jq query to execute on result")
	rootCmd.PersistentFlags().StringVarP(&optSession, "session", "S", "", "Session ID for private projects")

	debug := rootCmd.PersistentFlags().Bool("debug", false, "Turn on debug information")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if *debug {
			http.DefaultTransport = loghttp.DefaultTransport
		}
	}
}
