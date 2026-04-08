package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fumo-cli/fumo-command-line-interface/internal/fumo"
)

var BuildDate = "unknown"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of fumo CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("fumo CLI %s (built %s)\n", fumo.Version, BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
