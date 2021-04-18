package ui

import (
	"chat/shared"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func InitUI(receiver chan shared.TerminalData, sending chan string) {
	p := tea.NewProgram(initialModel(receiver, sending))
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
