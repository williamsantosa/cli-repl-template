package cmd

import (
	"github.com/spf13/cobra"

	"github.com/williamsantosa/cli-repl-template/internal/app"
)

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Interactive REPL with animated art",
	Long:  "Starts an interactive command prompt with the art animating at the top. Type commands below while it plays.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.RunREPL()
	},
}

func init() {
	rootCmd.AddCommand(replCmd)
}
