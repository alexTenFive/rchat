package ui

import (
	"chat/shared"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	model struct {
		sub          chan shared.TerminalData
		sending      chan string
		quitting     bool
		messageInput textinput.Model
		messages     []string // chat messages
	}
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

func initialModel(ch chan shared.TerminalData, sending chan string) model {
	name := textinput.NewModel()
	name.Placeholder = "start typing message..."
	name.Focus()
	name.PromptStyle = focusedStyle
	name.TextStyle = focusedStyle
	name.CharLimit = 32
	return model{ch, sending, false, name, []string{}}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, waitForActivity(m.sub))
}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan shared.TerminalData) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+C":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			m.sending <- m.messageInput.View()
			m.messageInput.SetValue("")
		}
	case shared.TerminalData: // record external activity
		m.messages = append(m.messages, msg.Message)
	}
	m.messageInput, cmd = m.messageInput.Update(msg)
	return m, tea.Batch(cmd, waitForActivity(m.sub))
}

func (m model) View() string {
	// The header
	s := "Welcome to the rchat. The most excellent chat in the world\n\n"

	for _, msg := range m.messages {
		s += msg + "\n"
	}

	s += m.messageInput.View() + "\n"
	// Send the UI for rendering
	return s
}
