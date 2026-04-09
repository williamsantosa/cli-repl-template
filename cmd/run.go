package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/williamsantosa/cli-repl-template/internal/app"
)

var runCmd = &cobra.Command{
	Use:   "run <command> [args...]",
	Short: "Run a REPL command non-interactively",
	Long:  "Executes a single REPL command and prints the result.\nUse 'run help' to see available REPL commands.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := strings.Join(args, " ")
		output := app.ExecuteCommand(input)
		if output != "" {
			fmt.Println(output)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
