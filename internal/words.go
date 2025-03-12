package words

import (
	"os"
	"path/filepath"
	"strings"
)

const DIR = "../assets"

func LoadWords(lang string) ([]string, error) {
	filePath := filepath.Join(DIR, lang+".txt")
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	words := strings.Split(string(fileData), "\n")

	return words, nil
}
