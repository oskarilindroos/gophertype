package words

import (
	_ "embed"
	"math/rand"
	"strings"
)

type WordManager struct {
	words []string
}

//go:embed assets/english.txt
var englishWords string

func NewWordManager() *WordManager {
	words := strings.Fields(englishWords)
	return &WordManager{
		words: words,
	}
}

// Returns all words
func (wm *WordManager) GetAllWords() []string {
	return wm.words
}

// Return n amount of random words
func (wm *WordManager) GetRandomWords(n int) []string {
	wordsLength := len(wm.words)

	// Limit n
	if n > wordsLength {
		n = wordsLength
	}

	randomWords := []string{}
	for range n {
		idx := rand.Intn(wordsLength)
		randomWords = append(randomWords, strings.ToLower(wm.words[idx]))
	}

	return randomWords
}
