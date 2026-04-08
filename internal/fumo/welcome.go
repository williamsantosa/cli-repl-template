package fumo

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fumo-cli/fumo-command-line-interface/internal/config"
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
	accent := lipgloss.Color("205")
	bright := lipgloss.Color("255")
	text := lipgloss.Color("252")
	subtle := lipgloss.Color("240")

	greetStyle := lipgloss.NewStyle().Bold(true).Foreground(bright)
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(accent)
	tipStyle := lipgloss.NewStyle().Foreground(text)
	boldTip := lipgloss.NewStyle().Bold(true).Foreground(bright)
	dimStyle := lipgloss.NewStyle().Foreground(subtle)

	greeting := greetStyle.Render("Welcome back, " + getUsername() + "!")

	// Left column: greeting + art
	left := "  " + greeting + "\n\n" + artFrame

	// Right column: tips, config, cwd
	source := config.C.Art.Source
	if source != "built-in" {
		source = filepath.Base(source)
	}

	right := headerStyle.Render("Tips for getting started") + "\n" +
		tipStyle.Render("Type ") + boldTip.Render("help") + tipStyle.Render(" to see commands") + "\n" +
		tipStyle.Render("Type ") + boldTip.Render("exit") + tipStyle.Render(" to leave the REPL") + "\n\n" +
		headerStyle.Render("Config") + "\n" +
		dimStyle.Render("Art: ") + tipStyle.Render(source) + "\n" +
		dimStyle.Render("Width: ") + tipStyle.Render(fmt.Sprint(config.C.Art.Width))

	if cwd := getCwd(); cwd != "" {
		right += "\n\n" + dimStyle.Render(cwd)
	}

	// Equalise column heights
	lh := lipgloss.Height(left)
	rh := lipgloss.Height(right)
	h := max(lh, rh)
	if lh < h {
		left += strings.Repeat("\n", h-lh)
	}
	if rh < h {
		right += strings.Repeat("\n", h-rh)
	}

	// Vertical divider
	divStyle := lipgloss.NewStyle().Foreground(accent)
	divParts := make([]string, h)
	for i := range divParts {
		divParts[i] = divStyle.Render("│")
	}
	divider := strings.Join(divParts, "\n")

	inner := lipgloss.JoinHorizontal(lipgloss.Top, left+"  ", divider, "  "+right)

	box := wrapInBox(inner, " fumo cli "+Version+" ", accent)
	hint := dimStyle.Render("  Press any key to continue...")

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

	// Top border: ╭─ title ────────╮
	titleRendered := tc.Render(title)
	titleW := lipgloss.Width(titleRendered)
	dashes := innerW - 1 - titleW
	if dashes < 0 {
		dashes = 0
	}
	top := bc.Render("╭─") + titleRendered + bc.Render(strings.Repeat("─", dashes)+"╮")

	// Blank line after top border for breathing room
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
