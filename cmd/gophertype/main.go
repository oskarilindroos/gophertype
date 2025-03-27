package main

import (
	"fmt"
	"gophertype/internal/words"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	containerStyle = lipgloss.NewStyle().
			Padding(2).
			BorderStyle(lipgloss.RoundedBorder())

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("yellow")).
			Background(lipgloss.Color("yellow")).
			Underline(true)

	correctWordStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("green"))

	incorrectWordStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("red"))
)

type model struct {
	width        int
	height       int
	userInput    string
	wordsManager *words.WordManager
	cursor       int
	correct      int
	errors       int
	wpm          int
}

func initialModel() model {
	wordsManager := words.NewWordManager()

	return model{
		wordsManager: wordsManager,
		cursor:       0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			return m, tea.Quit
		case "backspace":
			if m.cursor > 0 {
				m.cursor--
				m.userInput = m.userInput[:m.cursor]
			}
		default:
			m.cursor++
			m.userInput += msg.String()
		}
	// Update window width and height on resize
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	var output string

	// userInputWithCursor := m.userInput[:m.cursor] + cursorChar + m.userInput[m.cursor:]
	cursorChar := cursorStyle.Render(" ")
	userInputWithCursor := m.userInput[:m.cursor] + cursorChar + m.userInput[m.cursor:]

	output += containerStyle.
		Width(int(float64(m.width) * 0.5)). // Takes up half the width
		Render(fmt.Sprintf("%s", userInputWithCursor))

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, output)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
