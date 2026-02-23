package ui

import (
	"fmt"
	"strings"

	"github.com/TeamDei/tandem/internal/ui/layout"
	"github.com/TeamDei/tandem/internal/ui/styles"

	"github.com/charmbracelet/lipgloss"
)

// maxGlosses is the maximum number of glosses shown per part-of-speech group.
const maxGlosses = 5

// View dispatches to the correct view function for the current state.
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing…"
	}
	switch m.state {
	case StateReading:
		return m.viewReading()
	case StateLoading, StateLoadingDef:
		return m.viewLoading()
	case StateDetail:
		return m.viewDetail()
	case StateDefinition:
		return m.viewDefinition()
	case StateQuit:
		return m.viewQuit()
	default:
		panic(fmt.Sprintf("ui: unhandled AppState %d", m.state))
	}
}

func (m Model) viewReading() string {
	colW := layout.TextColWidth(m.width)
	th := layout.TextHeight(m.height)
	cw := layout.ContentWidth(m.width)

	l := layout.ComputeLayout(m.words, colW)
	numLines := len(l.LineStart)

	// Group word indices by their wrapped line.
	lineWords := make([][]int, numLines)
	for i, li := range l.LineOf {
		lineWords[li] = append(lineWords[li], i)
	}

	endLine := m.scrollOffset + th
	if endLine > numLines {
		endLine = numLines
	}

	var rows []string
	for li := m.scrollOffset; li < endLine; li++ {
		lineNum := styles.LineNum.Render(fmt.Sprintf("%4d", l.SourceLineOf[li]))
		sep := styles.Sep.Render(" │ ")
		var parts []string
		for _, wi := range lineWords[li] {
			if m.words[wi] == layout.Newline {
				continue // sentinel: not a visible word
			}
			if wi == m.cursor {
				parts = append(parts, styles.Cursor.Render(m.words[wi]))
			} else {
				parts = append(parts, styles.Word.Render(m.words[wi]))
			}
		}
		rows = append(rows, lineNum+sep+strings.Join(parts, " "))
	}

	// Pad to fill the panel height so the border doesn't shrink.
	for len(rows) < th {
		rows = append(rows, "")
	}

	body := lipgloss.NewStyle().Width(cw).Render(strings.Join(rows, "\n"))
	status := styles.Status.Render("  ←/h · →/l · ↑/k · ↓/j · g/G · enter analyze · tab define · q quit")

	return styles.Main.Width(layout.MainStyleWidth(m.width)).Render(
		styles.Title.Render("tandem") + "\n" + body + "\n" + status,
	)
}

func (m Model) viewLoading() string {
	content := fmt.Sprintf("\n %s  Looking up %q…\n", m.spinner.View(), m.selectedWord)
	panel := styles.Panel.Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, panel)
}

func (m Model) viewDetail() string {
	word := styles.PanelTitle.Render(m.selectedWord)
	sep := styles.Dim.Render(strings.Repeat("─", 44))

	var lines []string
	lines = append(lines, word, sep)

	if m.err != nil {
		lines = append(lines, styles.Analysis.Render("  Error: "+m.err.Error()))
	} else if len(m.response.Analyses) == 0 {
		lines = append(lines, styles.Analysis.Render("  No analyses found."))
	} else {
		for _, a := range m.response.Analyses {
			s := a.String()
			if len(s) > 0 && s[0] == '(' {
				if end := strings.Index(s, ")"); end > 0 {
					lines = append(lines, "  "+styles.POS.Render(s[:end+1])+styles.Analysis.Render(s[end+1:]))
					continue
				}
			}
			lines = append(lines, "  "+styles.Analysis.Render(s))
		}
	}

	lines = append(lines, "", styles.Status.Render("  esc / q  ·  back to reader"))

	panel := styles.Panel.Width(m.width / 2).Render(strings.Join(lines, "\n"))
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, panel)
}

func (m Model) viewDefinition() string {
	title := m.dictBaseWord
	if title == "" {
		title = m.selectedWord
	}
	word := styles.PanelTitle.Render(title)

	// If the inflected form differs from the lemma, show the source word.
	subtitle := ""
	if m.dictBaseWord != "" && m.dictBaseWord != m.selectedWord {
		subtitle = "\n  " + styles.Dim.Render("← "+m.selectedWord)
	}
	sep := styles.Dim.Render(strings.Repeat("─", 44))

	var lines []string
	lines = append(lines, word+subtitle, sep)

	switch {
	case m.err != nil:
		lines = append(lines, styles.Analysis.Render("  Error: "+m.err.Error()))
	case len(m.definitions) == 0:
		lines = append(lines, styles.Analysis.Render("  No definitions found."))
	default:
		for i, d := range m.definitions {
			if i > 0 {
				lines = append(lines, "") // blank line between POS groups
			}
			lines = append(lines, "  "+styles.POS.Render(d.PartOfSpeech))
			limit := len(d.Glosses)
			if limit > maxGlosses {
				limit = maxGlosses
			}
			for j, g := range d.Glosses[:limit] {
				lines = append(lines, fmt.Sprintf("  %s", styles.Analysis.Render(fmt.Sprintf("%d. %s", j+1, g))))
			}
		}
	}

	lines = append(lines, "", styles.Status.Render("  esc / q  ·  back to reader"))

	panel := styles.Panel.Width(m.width / 2).Render(strings.Join(lines, "\n"))
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, panel)
}

func (m Model) viewQuit() string {
	content := styles.PanelTitle.Render("Quit tandem?") + "\n\n" +
		styles.Word.Render("  y / enter  ·  quit") + "\n" +
		styles.Word.Render("  n / esc    ·  cancel")
	panel := styles.Quit.Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, panel)
}
