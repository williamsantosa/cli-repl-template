package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

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
	Long:  banner + "\n\nA configurable CLI template featuring a loading screen with art.\nAdd your own subcommands to extend it.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ./config.yaml or $HOME/.cli-repl/config.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
}

func initConfig() {
	if err := config.Load(cfgFile); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
	}
}
