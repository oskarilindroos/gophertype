package main

import (
	"flag"
	"fmt"
	"gophertype/internal/style"
	"gophertype/internal/words"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	charsPerWord     = 5   // Standard number of characters that count as one word
	minTimeElapsed   = 0.5 // Minimum time elapsed in minutes to prevent division by zero
	initialWordCount = 20  // Number of words to start out with
	minWordsLeft     = 8   // Minimum number of words left before generating more
	wordsToGenerate  = 10  // Number of words to generate when needed
)

type tickMsg time.Time

type model struct {
	windowWidth  int
	windowHeight int
	userInput    string
	wordsManager *words.WordManager // Word manager dependency for generating words
	targetWords  []string
	cursor       int       // Cursor position in the input
	correct      int       // Number of correctly typed words
	errors       int       // Number of incorrectly typed words
	finished     bool      // Indicates if the test is finished
	currentWord  int       // Index of the current word
	wordCursor   int       // Cursor position within the current word
	timeLeft     int       // Time left in seconds
	started      bool      // Indicates if the test has started
	startTime    time.Time // Start time of the test
	endTime      time.Time // End time of the test
}

func initialModel(testDuration int) model {
	wordsManager := words.NewWordManager()
	targetWords := wordsManager.GetRandomWords(initialWordCount)

	return model{
		wordsManager: wordsManager,
		targetWords:  targetWords,
		cursor:       0,
		finished:     false,
		currentWord:  0,
		wordCursor:   0,
		timeLeft:     testDuration,
		started:      false,
		startTime:    time.Time{},
		endTime:      time.Time{},
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("gophertype")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle quit signals first
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			return m, tea.Quit
		}

		// Don't start the test if it's already finished
		if !m.started && !m.finished {
			m.started = true
			m.startTime = time.Now()
			return m, tea.Batch(
				tea.Tick(time.Second, func(t time.Time) tea.Msg {
					return tickMsg(t)
				}),
				func() tea.Msg {
					return msg
				},
			)
		}

		// Don't process any other inputs if finished
		if m.finished {
			return m, nil
		}

		switch msg.String() {
		case "backspace":
			if m.wordCursor > 0 {
				m.wordCursor--
				m.userInput = m.userInput[:len(m.userInput)-1]
			}
		case " ":
			// Only allow pressing space to move to next word if at least one character has been typed
			if m.wordCursor > 0 {
				if m.currentWord < len(m.targetWords)-1 {
					// Check if current word is correct
					currentTarget := m.targetWords[m.currentWord]
					typedWord := m.userInput[len(m.userInput)-m.wordCursor:]
					if typedWord == currentTarget {
						m.correct++
					} else {
						m.errors++
					}
					m.currentWord++
					m.wordCursor = 0
					m.userInput += " "

					// Generate more words
					if len(m.targetWords)-m.currentWord < minWordsLeft {
						moreWords := m.wordsManager.GetRandomWords(wordsToGenerate)
						m.targetWords = append(m.targetWords, moreWords...)
					}
				} else {
					// Last word
					currentTarget := m.targetWords[m.currentWord]
					typedWord := m.userInput[len(m.userInput)-m.wordCursor:]
					if typedWord == currentTarget {
						m.correct++
					} else {
						m.errors++
					}
					m.finished = true
					m.endTime = time.Now()
				}
			}
		// Handle all other key presses
		default:
			if m.currentWord < len(m.targetWords) {
				currentTarget := m.targetWords[m.currentWord]
				if m.wordCursor < len(currentTarget) {
					m.userInput += msg.String()
					m.wordCursor++
				}
			}
		}
	// Update window size on resize
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		return m, nil
	case tickMsg:
		if m.started && !m.finished {
			m.timeLeft--
			if m.timeLeft <= 0 {
				m.finished = true
				m.endTime = time.Now() // Stop the timer
				return m, nil
			}
			return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		}
	}

	return m, nil
}

func calculateWPM(userInput string, timeElapsed time.Duration, errors int) (grossWPM int, netWPM int) {
	minutes := timeElapsed.Minutes()
	if minutes == 0 {
		minutes = minTimeElapsed
	}

	gross := float64(len(userInput)) / charsPerWord / minutes
	net := max((float64(len(userInput))/charsPerWord-float64(errors))/minutes, 0)

	return int(gross), int(net)
}

func calculateAccuracy(correct, errors int) float64 {
	if correct+errors == 0 {
		return 100.0
	}
	return float64(correct) / float64(correct+errors) * 100
}

func formatWordDisplay(targetWord, typedPart, currentChar, remainingPart string) string {
	var wordDisplay string

	// Add the typed part with correct/incorrect styling
	if len(typedPart) > 0 {
		for j, char := range typedPart {
			if j < len(targetWord) && string(char) == string(targetWord[j]) {
				wordDisplay += style.TypedStyle.Render(string(char))
			} else {
				wordDisplay += style.IncorrectStyle.Render(string(char))
			}
		}
	}

	// Add the current character with cursor
	if len(currentChar) > 0 {
		wordDisplay += style.CursorStyle.Render(currentChar)
	}

	// Add the remaining part
	wordDisplay += style.PlaceholderStyle.Render(remainingPart)
	return wordDisplay
}

func (m model) View() string {
	var output string

	if m.finished {
		// Calculate statistics
		timeElapsed := m.endTime.Sub(m.startTime)
		grossWPM, netWPM := calculateWPM(m.userInput, timeElapsed, m.errors)
		accuracy := calculateAccuracy(m.correct, m.errors)

		// Format and display results
		stats := fmt.Sprintf(
			"%s\n\n"+
				"WPM\n%s\n\n"+
				"Raw WPM\n%s\n\n"+
				"Accuracy\n%s\n\n"+
				"Correct\n%s\n\n"+
				"Errors\n%s",
			style.ResultsTitleStyle.Render("Results"),
			style.ResultsValueStyle.Render(fmt.Sprintf("%d", netWPM)),
			style.ResultsValueStyle.Render(fmt.Sprintf("%d", grossWPM)),
			style.ResultsValueStyle.Render(fmt.Sprintf("%.1f%%", accuracy)),
			style.ResultsValueStyle.Render(fmt.Sprintf("%d", m.correct)),
			style.ResultsValueStyle.Render(fmt.Sprintf("%d", m.errors)))

		output += style.ContainerStyle.
			Width(int(float64(m.windowWidth) * 0.5)).
			Align(lipgloss.Center).
			Render(stats)
	} else {
		var displayWords []string

		for i, targetWord := range m.targetWords {
			if i < m.currentWord {
				// Past words
				typedWord := strings.Fields(m.userInput)[i]
				if typedWord == targetWord {
					displayWords = append(displayWords, style.CorrectStyle.Render(targetWord))
				} else {
					displayWords = append(displayWords, style.IncorrectStyle.Render(targetWord))
				}
			} else if i == m.currentWord {
				// Current word
				typedPart := m.userInput[len(m.userInput)-m.wordCursor:]
				remainingPart := targetWord[m.wordCursor:]

				// Split the remaining part to get the current character
				var currentChar string
				if len(remainingPart) > 0 {
					currentChar = string(remainingPart[0])
					remainingPart = remainingPart[1:]
				}

				wordDisplay := formatWordDisplay(targetWord, typedPart, currentChar, remainingPart)
				displayWords = append(displayWords, wordDisplay)
			} else {
				// Future words
				displayWords = append(displayWords, style.PlaceholderStyle.Render(targetWord))
			}
		}

		// Create the typing area with timer
		typingArea := strings.Join(displayWords, " ")
		timer := style.TimerStyle.Render(fmt.Sprintf("%d", m.timeLeft))
		typingArea = timer + "\n" + typingArea

		output += style.ContainerStyle.
			Width(int(float64(m.windowWidth) * 0.5)).
			Render(typingArea)
	}

	return lipgloss.Place(m.windowWidth, m.windowHeight, lipgloss.Center, lipgloss.Center, output)
}

func main() {
	// Command-line flags
	testDuration := flag.Int("time", 30, "Duration of the typing test in seconds")
	flag.Parse()

	if *testDuration <= 0 {
		fmt.Println("Error: Test duration must be greater than 0 seconds")
		os.Exit(1)
	}

	p := tea.NewProgram(
		initialModel(*testDuration),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
