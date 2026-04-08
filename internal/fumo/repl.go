package fumo

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CommandHandler processes a command string and returns output.
// Register your own handlers via RegisterCommand before calling RunREPL.
type CommandHandler func(args []string) string

var commands map[string]CommandHandler

func init() {
	commands = map[string]CommandHandler{
		"help": func(args []string) string {
			var sb strings.Builder
			sb.WriteString("Available commands:\n")
			for name := range commands {
				sb.WriteString("  " + name + "\n")
			}
			return sb.String()
		},
		"echo": func(args []string) string {
			return strings.Join(args, " ")
		},
	}
}

// RegisterCommand adds a named command handler to the REPL.
func RegisterCommand(name string, handler CommandHandler) {
	commands[name] = handler
}

type replModel struct {
	input    textinput.Model
	quitting bool

	promptStyle lipgloss.Style
	outputStyle lipgloss.Style
	dimStyle    lipgloss.Style
}

func newREPLModel() replModel {
	ti := textinput.New()
	ti.Placeholder = "type a command..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60
	ti.Prompt = "fumo> "
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	return replModel{
		input:       ti,
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
			line := m.promptStyle.Render("fumo> ") + m.outputStyle.Render(input)
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

func (m replModel) executeCommand(input string) string {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return ""
	}
	name := strings.ToLower(parts[0])
	args := parts[1:]

	if handler, ok := commands[name]; ok {
		return handler(args)
	}
	return fmt.Sprintf("unknown command: %s (type 'help' for available commands)", name)
}

func (m replModel) View() string {
	if m.quitting {
		return ""
	}
	return m.input.View() + "\n" +
		m.dimStyle.Render("  exit/quit to leave")
}

// RunREPL plays the fumo animation once, then starts the interactive REPL
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
