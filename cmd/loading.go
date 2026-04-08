package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/williamsantosa/cli-repl-template/internal/app"
)

var loadingDuration int

var loadingCmd = &cobra.Command{
	Use:   "loading",
	Short: "Demo the loading animation",
	Long:  "Shows the art with an animated spinner for a configurable duration.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.RunLoader("Loading something awesome...", func() error {
			time.Sleep(time.Duration(loadingDuration) * time.Second)
			return nil
		})
	},
}

func init() {
	loadingCmd.Flags().IntVarP(&loadingDuration, "duration", "d", 3, "how many seconds to show the loader")
	rootCmd.AddCommand(loadingCmd)
}
