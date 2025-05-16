# gophertype

A fast, terminal-based typing test inspired by MonkeyType, built with Go, Bubble Tea, and Lipgloss.

## Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/oskarilindroos/gophertype.git
   cd gophertype
   ```
2. **Build the project:**
   ```sh
   go build -o gophertype ./cmd/gophertype
   ```

## Usage

Run the typing test with the default 30-second timer:
```sh
./gophertype
```

You can customize the test duration (in seconds):
```sh
./gophertype -time 60
```

## Customization
- **Word List:** Edit `internal/words/assets/english.txt` to use your own words.
- **Test Duration:** Use the `-time` flag to set the timer.

## Dependencies
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)

## License
MIT License
