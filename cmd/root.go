package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/fumo-cli/fumo-command-line-interface/internal/config"
)

var cfgFile string

var banner = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	Render("  ╒══════════════════╕\n  │   fumo cli  ◕ᴗ◕  │\n  ╘══════════════════╛")

var rootCmd = &cobra.Command{
	Use:   "fumo",
	Short: "Fumo CLI — a command-line interface with fumo flair",
	Long:  banner + "\n\nA configurable CLI template featuring a fumo loading screen.\nAdd your own subcommands to extend it.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ./fumo.yaml or $HOME/.fumo/fumo.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
}

func initConfig() {
	if err := config.Load(cfgFile); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
	}
}
