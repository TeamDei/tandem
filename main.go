package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/servusdei2018/tandem/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	file string
)

func main() {
	flag.StringVar(&file, "file", "", "path to file")
	flag.Parse()

	if file == "" {
		fmt.Fprintln(os.Stderr, "error: --file is required. Usage: tandem --file <path>")
		os.Exit(1)
	}

	raw, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	words := parseWords(raw)
	if len(words) == 0 {
		fmt.Fprintln(os.Stderr, "error: file is empty or contains no words")
		os.Exit(1)
	}

	p := tea.NewProgram(ui.New(words), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// parseWords splits raw file bytes into a word list, inserting "\n" sentinel
// tokens to represent hard line-breaks from the source. Consecutive blank lines
// each produce a single sentinel so the reader gets one blank visual line.
func parseWords(raw []byte) []string {
	lines := strings.Split(strings.ReplaceAll(string(raw), "\r", ""), "\n")
	var words []string
	prevHadWords := false
	for _, line := range lines {
		lineWords := strings.Fields(line)
		if len(lineWords) == 0 {
			if prevHadWords {
				words = append(words, "\n") // blank visual line
			}
			continue
		}
		if prevHadWords {
			words = append(words, "\n") // forced line break
		}
		words = append(words, lineWords...)
		prevHadWords = true
	}
	return words
}
