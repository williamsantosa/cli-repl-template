package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/williamsantosa/cli-repl-template/internal/app"
	"github.com/williamsantosa/cli-repl-template/internal/config"
)

var cfgFile string

var banner = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	Render("  ╒══════════════════╕\n  │    cli-repl      │\n  ╘══════════════════╛")

var rootCmd = &cobra.Command{
	Use:   "cli-repl",
	Short: "A configurable CLI REPL template",
	Long: banner + "\n\n" +
		"A configurable CLI template featuring a loading screen with art.\n" +
		"Run with no subcommand to start the interactive REPL (same as the 'repl' command).\n" +
		"Add your own subcommands to extend it.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.RunREPL()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ./config.yaml or $HOME/.cli-repl/config.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output (debug log to stderr)")
	rootCmd.PersistentPreRunE = persistentSetup
}

func configureLogging(verbose bool) {
	var h slog.Handler
	if verbose {
		h = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		h = slog.DiscardHandler
	}
	slog.SetDefault(slog.New(h))
}

func persistentSetup(cmd *cobra.Command, args []string) error {
	verbose, err := cmd.Root().PersistentFlags().GetBool("verbose")
	if err != nil {
		return err
	}
	configureLogging(verbose)
	if err := config.Load(cfgFile); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
	}
	return nil
}
