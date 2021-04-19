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
		quitting     chan bool
		theight      uint
		messageInput textinput.Model
		messages     []string // chat messages
	}
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

func initialModel(ch chan shared.TerminalData, sending chan string, quitting chan bool) model {
	name := textinput.NewModel()
	name.Placeholder = "start typing message..."
	name.Focus()
	name.PromptStyle = focusedStyle
	name.TextStyle = focusedStyle
	name.CharLimit = 256
	theight := getHeight()
	return model{ch, sending, quitting, theight, name, []string{}}
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
		case "ctrl+c", "ctrl+C", "esc":
			defer func() {
				m.messageInput.Reset()
				m.quitting <- true
			}()
			return m, tea.Quit
		case "enter":
			m.sending <- m.messageInput.View()
			m.messageInput.Reset()

		}
	case shared.TerminalData: // record external activity
		m.messages = append(m.messages, msg.Message)
	}
	m.messageInput, cmd = m.messageInput.Update(msg)
	return m, tea.Batch(cmd, waitForActivity(m.sub))
}

func (m model) View() string {
	// The header
	s := "\nWelcome to the rchat. The most excellent chat in the world\n"
	emptyLines := 0
	height := m.theight - 2
	if uint(len(m.messages)) < height {
		emptyLines = int(height) - len(m.messages)
	}

	for i := 0; i < emptyLines; i++ {
		s += "\n"
	}

	first, last := len(m.messages)-int(height), len(m.messages)
	if first < 0 {
		first = 0
	}
	for _, msg := range m.messages[first:last] {
		s += msg + "\n"
	}

	s += m.messageInput.View() + "\n"

	// Send the UI for rendering
	return s
}
