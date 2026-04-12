package app

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/williamsantosa/cli-repl-template/internal/config"
)

// Command holds a handler and its description for the REPL.
type Command struct {
	Handler     func(args []string) string
	Description string
}

var commands map[string]Command

func init() {
	commands = map[string]Command{
		"help": {
			Description: "Show available commands, or 'help <command>' for details",
			Handler: func(args []string) string {
				if len(args) > 0 {
					name := strings.ToLower(args[0])
					if cmd, ok := commands[name]; ok {
						return name + " — " + cmd.Description
					}
					return fmt.Sprintf("unknown command: %s", name)
				}
				var sb strings.Builder
				sb.WriteString("Available commands:\n")
				names := make([]string, 0, len(commands))
				for name := range commands {
					names = append(names, name)
				}
				slices.Sort(names)
				for _, name := range names {
					cmd := commands[name]
					sb.WriteString(fmt.Sprintf("  %-10s %s\n", name, cmd.Description))
				}
				return strings.TrimRight(sb.String(), "\n")
			},
		},
		"echo": {
			Description: "Repeat the given text back",
			Handler: func(args []string) string {
				return strings.Join(args, " ")
			},
		},
	}
}

// RegisterCommand adds a named command handler to the REPL.
func RegisterCommand(name string, description string, handler func(args []string) string) {
	commands[name] = Command{Handler: handler, Description: description}
}

type replModel struct {
	input    textinput.Model
	quitting bool
	prompt   string

	promptStyle lipgloss.Style
	outputStyle lipgloss.Style
	dimStyle    lipgloss.Style
}

func newREPLModel() replModel {
	prompt := config.C.Name + "> "

	ti := textinput.New()
	ti.Placeholder = "type a command..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60
	ti.Prompt = prompt
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	return replModel{
		input:       ti,
		prompt:      prompt,
		promptStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true),
		outputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		dimStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
	}
}

func (m replModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m replModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			input := strings.TrimSpace(m.input.Value())
			m.input.SetValue("")
			if input == "" {
				return m, nil
			}
			if input == "exit" || input == "quit" {
				m.quitting = true
				return m, tea.Quit
			}
			output := m.executeCommand(input)
			if output == "\x00CLEAR" {
				return m, tea.ClearScreen
			}
			line := m.promptStyle.Render(m.prompt) + m.outputStyle.Render(input)
			if output != "" {
				line += "\n" + m.outputStyle.Render(output)
			}
			return m, tea.Println(line)
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// ExecuteCommand dispatches a single REPL command string and returns the output.
// This is used by both the interactive REPL and the non-interactive "run" subcommand.
func ExecuteCommand(input string) string {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return ""
	}
	name := strings.ToLower(parts[0])
	args := parts[1:]

	if cmd, ok := commands[name]; ok {
		return cmd.Handler(args)
	}
	return fmt.Sprintf("unknown command: %s (type 'help' for available commands)", name)
}

func (m replModel) executeCommand(input string) string {
	return ExecuteCommand(input)
}

func (m replModel) View() string {
	if m.quitting {
		return ""
	}
	return m.dimStyle.Render("──────────────────────────────────────────────────") + "\n" +
		m.input.View() + "\n" +
		m.dimStyle.Render("  type 'help' for commands · exit/quit to leave")
}

// RunREPL plays the animation once, then starts the interactive REPL
// with the prompt pinned at the bottom. Command output is printed into the
// terminal scrollback above the prompt.
func RunREPL() error {
	if err := RunAnimationUntilInput(); err != nil {
		return err
	}

	m := newREPLModel()
	_, err := tea.NewProgram(m).Run()
	return err
}
