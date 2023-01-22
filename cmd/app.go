package cmd

import (
	"log"

	"github.com/haikalvidya/go-article/internal/app"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "go-article is a REST API service article",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = app.Run(
			app.TypeHTTPServer,
			app.WithArgs(args),
		)

		return
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing the root cmd: %v.", err)
	}
}
