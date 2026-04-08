package app

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/williamsantosa/cli-repl-template/internal/config"
)

// Version can be set at build time via ldflags.
var Version = "dev"

func getUsername() string {
	u, err := user.Current()
	if err != nil {
		return "friend"
	}
	name := u.Username
	if i := strings.LastIndexByte(name, '\\'); i >= 0 {
		name = name[i+1:]
	}
	return name
}

func getCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	home, err := os.UserHomeDir()
	if err == nil {
		sep := string(os.PathSeparator)
		if dir == home {
			return "~"
		}
		if strings.HasPrefix(dir, home+sep) {
			return "~" + sep + dir[len(home)+1:]
		}
	}
	return dir
}

func renderWelcomeView(artFrame string) string {
	cfg := config.C.Welcome
	accent := lipgloss.Color(cfg.AccentColor)
	bright := lipgloss.Color("255")
	text := lipgloss.Color("252")
	subtle := lipgloss.Color("240")

	greetStyle := lipgloss.NewStyle().Bold(true).Foreground(bright)
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(accent)
	tipStyle := lipgloss.NewStyle().Foreground(text)
	dimStyle := lipgloss.NewStyle().Foreground(subtle)

	greetText := strings.ReplaceAll(cfg.Greeting, "{user}", getUsername())
	greeting := greetStyle.Render(greetText)

	left := "  " + greeting + "\n\n" + artFrame

	var rightParts []string

	if cfg.ShowTips && len(cfg.Tips) > 0 {
		section := headerStyle.Render(cfg.TipsTitle)
		for _, t := range cfg.Tips {
			section += "\n" + tipStyle.Render(t)
		}
		rightParts = append(rightParts, section)
	}

	if cfg.ShowConfig {
		source := config.C.Art.Source
		if source != "built-in" {
			source = filepath.Base(source)
		}
		section := headerStyle.Render("Config") + "\n" +
			dimStyle.Render("Art: ") + tipStyle.Render(source) + "\n" +
			dimStyle.Render("Width: ") + tipStyle.Render(fmt.Sprint(config.C.Art.Width))
		rightParts = append(rightParts, section)
	}

	if cfg.ShowCwd {
		if cwd := getCwd(); cwd != "" {
			rightParts = append(rightParts, dimStyle.Render(cwd))
		}
	}

	right := strings.Join(rightParts, "\n\n")

	lh := lipgloss.Height(left)
	rh := lipgloss.Height(right)
	h := max(lh, rh)
	if lh < h {
		left += strings.Repeat("\n", h-lh)
	}
	if rh < h {
		right += strings.Repeat("\n", h-rh)
	}

	divStyle := lipgloss.NewStyle().Foreground(accent)
	divParts := make([]string, h)
	for i := range divParts {
		divParts[i] = divStyle.Render("│")
	}
	divider := strings.Join(divParts, "\n")

	inner := lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", divider, "  ", right)

	box := wrapInBox(inner, " "+config.C.Name+" cli "+Version+" ", accent)
	hint := dimStyle.Render("  " + cfg.Hint)

	return box + "\n\n" + hint + "\n"
}

// wrapInBox draws a rounded-corner box around content with a title
// embedded in the top border line.
func wrapInBox(content string, title string, borderColor lipgloss.Color) string {
	bc := lipgloss.NewStyle().Foreground(borderColor)
	tc := lipgloss.NewStyle().Bold(true).Foreground(borderColor)

	lines := strings.Split(content, "\n")

	maxW := 0
	for _, line := range lines {
		if w := lipgloss.Width(line); w > maxW {
			maxW = w
		}
	}

	innerW := maxW + 4

	titleRendered := tc.Render(title)
	titleW := lipgloss.Width(titleRendered)
	dashes := innerW - 1 - titleW
	if dashes < 0 {
		dashes = 0
	}
	top := bc.Render("╭─") + titleRendered + bc.Render(strings.Repeat("─", dashes)+"╮")

	blankPad := strings.Repeat(" ", innerW)
	spacer := bc.Render("│") + blankPad + bc.Render("│")

	var body strings.Builder
	for _, line := range lines {
		pad := maxW - lipgloss.Width(line)
		if pad < 0 {
			pad = 0
		}
		body.WriteString(bc.Render("│") + "  " + line + strings.Repeat(" ", pad) + "  " + bc.Render("│") + "\n")
	}

	bottom := bc.Render("╰" + strings.Repeat("─", innerW) + "╯")

	return top + "\n" + spacer + "\n" + body.String() + spacer + "\n" + bottom
}
