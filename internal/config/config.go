package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type ArtConfig struct {
	Source      string  `mapstructure:"source"`
	Width       int     `mapstructure:"width"`
	CellRatio   float64 `mapstructure:"cell_ratio"`
	Border      bool    `mapstructure:"border"`
	BorderColor string  `mapstructure:"border_color"`
}

type LoaderConfig struct {
	Spinner      string `mapstructure:"spinner"`
	SpinnerColor string `mapstructure:"spinner_color"`
	SpeedMs      int    `mapstructure:"speed_ms"`
	MessageColor string `mapstructure:"message_color"`
}

type WelcomeConfig struct {
	Greeting    string   `mapstructure:"greeting"`
	Hint        string   `mapstructure:"hint"`
	AccentColor string   `mapstructure:"accent_color"`
	ShowTips    bool     `mapstructure:"show_tips"`
	TipsTitle   string   `mapstructure:"tips_title"`
	Tips        []string `mapstructure:"tips"`
	ShowConfig  bool     `mapstructure:"show_config"`
	ShowCwd     bool     `mapstructure:"show_cwd"`
}

type Config struct {
	Art     ArtConfig     `mapstructure:"art"`
	Loader  LoaderConfig  `mapstructure:"loader"`
	Welcome WelcomeConfig `mapstructure:"welcome"`
}

var C Config

func setDefaults() {
	viper.SetDefault("art.source", "built-in")
	viper.SetDefault("art.width", 40)
	viper.SetDefault("art.cell_ratio", 0.45)
	viper.SetDefault("art.border", true)
	viper.SetDefault("art.border_color", "63")

	viper.SetDefault("loader.spinner", "dots")
	viper.SetDefault("loader.spinner_color", "205")
	viper.SetDefault("loader.speed_ms", 100)
	viper.SetDefault("loader.message_color", "252")

	viper.SetDefault("welcome.greeting", "Welcome back, {user}!")
	viper.SetDefault("welcome.hint", "Press any key to continue...")
	viper.SetDefault("welcome.accent_color", "205")
	viper.SetDefault("welcome.show_tips", true)
	viper.SetDefault("welcome.tips_title", "Tips for getting started")
	viper.SetDefault("welcome.tips", []string{
		"Type 'help' to see commands",
		"Type 'exit' to leave the REPL",
	})
	viper.SetDefault("welcome.show_config", true)
	viper.SetDefault("welcome.show_cwd", true)
}

// Load reads configuration from the given file path (if non-empty),
// or searches the working directory and $HOME/.fumo/ for fumo.yaml.
func Load(cfgFile string) error {
	setDefaults()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("fumo")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		if home, err := os.UserHomeDir(); err == nil {
			viper.AddConfigPath(filepath.Join(home, ".fumo"))
		}
	}

	// It's fine if no config file exists — defaults will be used.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return viper.Unmarshal(&C)
}
