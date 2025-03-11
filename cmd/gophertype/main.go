package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

var containerStyle = lipgloss.NewStyle().
	Padding(2).
	BorderStyle(lipgloss.RoundedBorder())

type model struct {
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if s := msg.String(); s == "ctrl+c" || s == "q" || s == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	var s string

	s += containerStyle.
		Width(int(float64(m.width) * 0.5)).
		Render(fmt.Sprintf("60\nway where second group more who group end house even question miss new own four up night number still make word any way where second group more who group end house even question miss new own four up night number still make word any"))

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, s)
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
