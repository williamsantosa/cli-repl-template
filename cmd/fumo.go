package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/fumo-cli/fumo-command-line-interface/internal/fumo"
)

var loadingDuration int

var fumoCmd = &cobra.Command{
	Use:   "loading",
	Short: "Demo the fumo loading animation",
	Long:  "Shows the fumo art with an animated spinner for a configurable duration.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fumo.RunLoader("Loading something awesome...", func() error {
			time.Sleep(time.Duration(loadingDuration) * time.Second)
			return nil
		})
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display the fumo art (animated GIFs will loop until q/ctrl+c)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fumo.RunAnimation()
	},
}

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Interactive REPL with animated fumo art",
	Long:  "Starts an interactive command prompt with the fumo art animating at the top. Type commands below while the fumo dances.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fumo.RunREPL()
	},
}

func init() {
	fumoCmd.Flags().IntVarP(&loadingDuration, "duration", "d", 3, "how many seconds to show the loader")
	rootCmd.AddCommand(fumoCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(replCmd)
}
