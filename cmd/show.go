package cmd

import (
	"github.com/spf13/cobra"

	"github.com/williamsantosa/cli-repl-template/internal/app"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display the art (animated GIFs will loop until q/ctrl+c)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.RunAnimation()
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
