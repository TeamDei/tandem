package ui

import (
	"github.com/servusdei2018/tandem/internal/latin"
	"github.com/servusdei2018/tandem/internal/ui/layout"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// clamp returns v constrained to [lo, hi].
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// skipSentinel advances (or retreats) the cursor past any sentinel tokens
// so that the cursor always rests on a real word.
func skipSentinel(words []string, cursor, delta int) int {
	for cursor >= 0 && cursor < len(words) && words[cursor] == layout.Newline {
		cursor += delta
	}
	return clamp(cursor, 0, len(words)-1)
}

// AppState represents the current screen of the application.
type AppState int

const (
	StateReading AppState = iota
	StateLoading
	StateDetail
	StateLoadingDef
	StateDefinition
	StateQuit
)

// APIResultMsg is the message delivered when the Perseus API call completes.
type APIResultMsg struct {
	Resp latin.Response
	Err  error
}

// FetchWord returns a tea.Cmd that queries the Perseus API for word.
func FetchWord(word string) tea.Cmd {
	return func() tea.Msg {
		resp, err := latin.Lookup(word)
		return APIResultMsg{Resp: resp, Err: err}
	}
}

// DictResultMsg is the message delivered when the combined Perseus+Wiktionary call completes.
type DictResultMsg struct {
	BaseWord string // the lemma that was sent to Wiktionary
	Defs     []latin.Definition
	Err      error
}

// FetchDefinition returns a tea.Cmd that queries Wiktionary for English definitions of word.
// If the word only contains inflected forms, Wiktionary automatically follows the link to the base lemma.
func FetchDefinition(word string) tea.Cmd {
	return func() tea.Msg {
		base, defs, err := latin.LookupDefinitions(word)
		return DictResultMsg{BaseWord: base, Defs: defs, Err: err}
	}
}

// Model is the root bubbletea model for the tandem application.
type Model struct {
	words         []string
	cursor        int
	width, height int
	scrollOffset  int
	state         AppState
	spinner       spinner.Model
	response      latin.Response
	definitions   []latin.Definition
	selectedWord  string
	dictBaseWord  string // lemma sent to Wiktionary
	err           error
}

// New constructs an initialised Model for the given word list.
func New(words []string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))
	return Model{words: words, spinner: s}
}

// Init satisfies tea.Model; no initial command is needed.
func (m Model) Init() tea.Cmd { return nil }

// withScroll adjusts scrollOffset so the cursor line is always visible.
func (m Model) withScroll() Model {
	if len(m.words) == 0 || m.width == 0 {
		return m
	}
	l := layout.ComputeLayout(m.words, layout.TextColWidth(m.width))
	m.cursor = clamp(m.cursor, 0, len(l.LineOf)-1)
	curLine := l.LineOf[m.cursor]
	th := layout.TextHeight(m.height)
	if curLine < m.scrollOffset {
		m.scrollOffset = curLine
	} else if curLine >= m.scrollOffset+th {
		m.scrollOffset = curLine - th + 1
	}
	return m
}

// Update handles all incoming messages and key events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m.withScroll(), nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case APIResultMsg:
		m.response, m.err = msg.Resp, msg.Err
		m.state = StateDetail
		return m, nil

	case DictResultMsg:
		m.definitions, m.dictBaseWord, m.err = msg.Defs, msg.BaseWord, msg.Err
		m.state = StateDefinition
		return m, nil

	case tea.KeyMsg:
		switch m.state {

		case StateReading:
			return m.updateReading(msg)

		case StateDetail, StateDefinition:
			switch msg.String() {
			case "q", "esc":
				m.state = StateReading
			}

		case StateQuit:
			switch msg.String() {
			case "y", "enter":
				return m, tea.Quit
			case "n", "esc", "q":
				m.state = StateReading
			}
		}
	}

	return m, nil
}

// updateReading handles key input in the reading state.
func (m Model) updateReading(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	l := layout.ComputeLayout(m.words, layout.TextColWidth(m.width))
	numLines := len(l.LineStart)

	switch msg.String() {
	case "q", "esc":
		m.state = StateQuit

	case "enter":
		if len(m.words) == 0 || m.words[m.cursor] == layout.Newline {
			break
		}
		m.selectedWord = m.words[m.cursor]
		m.state = StateLoading
		return m, tea.Batch(m.spinner.Tick, FetchWord(m.selectedWord))

	case "tab":
		if len(m.words) == 0 || m.words[m.cursor] == layout.Newline {
			break
		}
		m.selectedWord = m.words[m.cursor]
		m.state = StateLoadingDef
		return m, tea.Batch(m.spinner.Tick, FetchDefinition(m.selectedWord))

	case "left", "h":
		if m.cursor > 0 {
			m.cursor = skipSentinel(m.words, m.cursor-1, -1)
		}

	case "right", "l":
		if m.cursor < len(m.words)-1 {
			m.cursor = skipSentinel(m.words, m.cursor+1, +1)
		}

	case "up", "k":
		curLine := l.LineOf[m.cursor]
		if curLine > 0 {
			col := m.cursor - l.LineStart[curLine]
			prevStart := l.LineStart[curLine-1]
			prevEnd := l.LineStart[curLine] - 1
			target := prevStart + col
			if target > prevEnd {
				target = prevEnd
			}
			m.cursor = skipSentinel(m.words, target, -1)
		}

	case "down", "j":
		curLine := l.LineOf[m.cursor]
		if curLine < numLines-1 {
			col := m.cursor - l.LineStart[curLine]
			nextStart := l.LineStart[curLine+1]
			var nextEnd int
			if curLine+2 < numLines {
				nextEnd = l.LineStart[curLine+2] - 1
			} else {
				nextEnd = len(m.words) - 1
			}
			target := nextStart + col
			if target > nextEnd {
				target = nextEnd
			}
			m.cursor = skipSentinel(m.words, target, +1)
		}

	case "g", "home":
		m.cursor, m.scrollOffset = 0, 0
		return m, nil

	case "G", "end":
		m.cursor = skipSentinel(m.words, len(m.words)-1, -1)

	case "ctrl+f", "pgdown":
		curLine := l.LineOf[m.cursor]
		newLine := curLine + layout.TextHeight(m.height)
		if newLine >= numLines {
			newLine = numLines - 1
		}
		m.cursor = skipSentinel(m.words, l.LineStart[newLine], +1)

	case "ctrl+b", "pgup":
		curLine := l.LineOf[m.cursor]
		newLine := curLine - layout.TextHeight(m.height)
		if newLine < 0 {
			newLine = 0
		}
		m.cursor = skipSentinel(m.words, l.LineStart[newLine], +1)
	}

	return m.withScroll(), nil
}
