package ui

import (
	"time"

	"cykl/sequencer"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	// We don't need to refresh the ui as often as the sequencer.
	// It saves some cpu. Right now we run it at 30 fps.
	refreshFrequency = 33 * time.Millisecond
)

// tickMsg is a message that triggers ui rrefresh
type tickMsg time.Time

type mainModel struct {
	seq    *sequencer.Sequencer
	width  int
	height int
}

// New creates a new mainModel that hols the ui state. It takes a new sequencer.
// Check teh sequencer package.
func New(seq *sequencer.Sequencer) tea.Model {
	model := mainModel{
		seq: seq,
	}
	return model
}

func tick() tea.Cmd {
	return tea.Tick(refreshFrequency, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tick())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			m.seq.TogglePlay()
			return m, nil
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m mainModel) View() string {
	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderSeq(),
	)

	// Cleanup gibber
	cleanup := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - lipgloss.Height(mainView)).
		Render("")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainView,
		cleanup,
	)
}

func (m mainModel) renderSeq() string {
	trackUi := []string{}
	for _, track := range m.seq.Tracks {
		trackUi = append(trackUi,
			lipgloss.NewStyle().
				MarginBottom(1).
				Render(m.renderTrack(track)),
		)
	}
	return lipgloss.JoinVertical(lipgloss.Left, trackUi...)
}

func (m mainModel) renderTrack(track *sequencer.Track) string {
	lines := []string{}
	for row := 0; row < 4; row++ {
		steps := []string{}
		for col := 0; col < 4; col++ {
			step := col + row*4
			if step >= track.Steps {
				steps = append(steps, " ")
				break
			}
			if m.seq.Playing && track.CurrentStep() == step {
				steps = append(steps, "░░")
			} else {
				steps = append(steps, "██")
			}
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, steps...))
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
