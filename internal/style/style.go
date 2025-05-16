package style

import "github.com/charmbracelet/lipgloss"

var (
	ContainerStyle = lipgloss.NewStyle().
			Padding(2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"}).
			Bold(true)

	CursorStyle = lipgloss.NewStyle().
			Underline(true).
			Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"})

	PlaceholderStyle = lipgloss.NewStyle().
				Faint(true).
				Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"})

	TypedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "000", Dark: "fff"})

	CorrectStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "28", Dark: "34"}).
			Bold(true)

	IncorrectStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "88", Dark: "160"}).
			Bold(true).
			Strikethrough(true)

	TimerStyle = lipgloss.NewStyle().
			Bold(true).
			Italic(true).
			Foreground(lipgloss.AdaptiveColor{Light: "236", Dark: "248"})

	ResultsTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Underline(true).
				Foreground(lipgloss.AdaptiveColor{Light: "000", Dark: "fff"})

	ResultsValueStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.AdaptiveColor{Light: "000", Dark: "fff"})
)
