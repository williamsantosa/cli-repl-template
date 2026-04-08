package fumo

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/fumo-cli/fumo-command-line-interface/internal/config"
)

type taskDoneMsg struct{ err error }
type frameTickMsg struct{}

type loaderModel struct {
	spinner  spinner.Model
	frames   []Frame
	frameIdx int
	message  string
	msgStyle lipgloss.Style
	done     bool
	err      error
}

func spinnerByName(name string) spinner.Spinner {
	switch name {
	case "line":
		return spinner.Line
	case "minidot":
		return spinner.MiniDot
	case "jump":
		return spinner.Jump
	case "pulse":
		return spinner.Pulse
	case "points":
		return spinner.Points
	case "globe":
		return spinner.Globe
	case "moon":
		return spinner.Moon
	case "monkey":
		return spinner.Monkey
	case "meter":
		return spinner.Meter
	case "hamburger":
		return spinner.Hamburger
	default:
		return spinner.Dot
	}
}

func newLoaderModel(message string) loaderModel {
	cfg := config.C.Loader

	sp := spinner.New()
	sp.Spinner = spinnerByName(cfg.Spinner)
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.SpinnerColor))

	if cfg.SpeedMs > 0 {
		frames := sp.Spinner.Frames
		sp.Spinner = spinner.Spinner{
			Frames: frames,
			FPS:    time.Duration(cfg.SpeedMs) * time.Millisecond,
		}
	}

	msgStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.MessageColor))

	return loaderModel{
		spinner:  sp,
		frames:   RenderFrames(),
		frameIdx: 0,
		message:  message,
		msgStyle: msgStyle,
	}
}

func (m loaderModel) scheduleFrameTick() tea.Cmd {
	if len(m.frames) <= 1 {
		return nil
	}
	delay := m.frames[m.frameIdx].Delay
	if delay == 0 {
		delay = 100 * time.Millisecond
	}
	return tea.Tick(delay, func(time.Time) tea.Msg {
		return frameTickMsg{}
	})
}

func (m loaderModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.scheduleFrameTick())
}

func (m loaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			m.done = true
			return m, tea.Quit
		}
	case taskDoneMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	case frameTickMsg:
		m.frameIdx = (m.frameIdx + 1) % len(m.frames)
		return m, m.scheduleFrameTick()
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m loaderModel) View() string {
	if m.done {
		return ""
	}
	art := m.frames[m.frameIdx].Rendered
	return fmt.Sprintf(
		"%s\n\n  %s %s\n",
		art,
		m.spinner.View(),
		m.msgStyle.Render(m.message),
	)
}

// RunLoader displays the fumo art with an animated spinner while task executes.
// If the art source is an animated GIF, the frames cycle automatically.
func RunLoader(message string, task func() error) error {
	m := newLoaderModel(message)

	p := tea.NewProgram(m)

	go func() {
		err := task()
		p.Send(taskDoneMsg{err: err})
	}()

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("loader error: %w", err)
	}

	if fm, ok := finalModel.(loaderModel); ok && fm.err != nil {
		return fm.err
	}

	return nil
}

// RunAnimation displays the fumo art in a loop (no spinner, no task).
// For animated GIFs, frames cycle until the user presses q or ctrl+c.
// For static art, it just displays until dismissed.
func RunAnimation() error {
	frames := RenderFrames()
	if len(frames) <= 1 {
		fmt.Println(frames[0].Rendered)
		return nil
	}

	m := animModel{frames: frames}
	_, err := tea.NewProgram(m).Run()
	return err
}

type animModel struct {
	frames      []Frame
	frameIdx    int
	once        bool
	waitForKey  bool
	done        bool
}

func (m animModel) Init() tea.Cmd {
	delay := m.frames[0].Delay
	if delay == 0 {
		delay = 100 * time.Millisecond
	}
	return tea.Tick(delay, func(time.Time) tea.Msg {
		return frameTickMsg{}
	})
}

func (m animModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.waitForKey || msg.String() == "ctrl+c" || msg.String() == "q" {
			m.done = true
			return m, tea.Quit
		}
	case frameTickMsg:
		next := (m.frameIdx + 1) % len(m.frames)
		if m.once && next == 0 {
			m.done = true
			return m, tea.Quit
		}
		m.frameIdx = next
		delay := m.frames[m.frameIdx].Delay
		if delay == 0 {
			delay = 100 * time.Millisecond
		}
		return m, tea.Tick(delay, func(time.Time) tea.Msg {
			return frameTickMsg{}
		})
	}
	return m, nil
}

func (m animModel) View() string {
	if m.done {
		return ""
	}
	if m.waitForKey {
		return renderWelcomeView(m.frames[m.frameIdx].Rendered)
	}
	return m.frames[m.frameIdx].Rendered + "\n"
}

// RunAnimationUntilInput loops the animation until the user presses any key.
// Shows a "Press any key to continue..." hint below the art.
// For static art, it displays and waits for a key press.
func RunAnimationUntilInput() error {
	frames := RenderFrames()
	m := animModel{frames: frames, waitForKey: true}
	_, err := tea.NewProgram(m).Run()
	return err
}
