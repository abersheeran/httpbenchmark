package core

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type counter[T any] struct {
	sub               chan T // where we'll receive activity notifications
	responses         int    // how many responses we've received
	spinner           spinner.Model
	listenForActivity func(chan T) tea.Cmd
	waitForActivity   func(chan T) tea.Cmd
}

func (m counter[T]) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.listenForActivity(m.sub),
		m.waitForActivity(m.sub),
	)
}

func (m counter[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case T:
		m.responses++                      // record external activity
		return m, m.waitForActivity(m.sub) // wait for next event
	case tea.KeyMsg:
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m counter[T]) View() string {
	s := fmt.Sprintf("\n %s Requests sent: %d\n   Press any key to exit\n\n", m.spinner.View(), m.responses)
	return s
}
