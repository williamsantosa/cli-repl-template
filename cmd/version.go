package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/williamsantosa/cli-repl-template/internal/app"
)

var BuildDate = "unknown"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cli-repl %s (built %s)\n", app.Version, BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
