package layout

// GutterVisualWidth is the fixed terminal-column width of the line-number gutter.
// "NNNN │ " = 4 digits + 1 space + │ + 1 space = 7 visible columns.
const GutterVisualWidth = 7

// MainStyleWidth returns the value to pass to styleMain.Width() so that the
// bordered panel fills the terminal.
//
// styleMain has Padding(0,1) (+2 horizontal) and a border (+2 horizontal).
// So total rendered width = Width arg + 2; to fill termW: Width arg = termW - 2.
func MainStyleWidth(termW int) int { return termW - 2 }

// ContentWidth is the number of usable characters inside the border and padding.
func ContentWidth(termW int) int { return termW - 4 }

// TextColWidth is the usable column width available for word text (gutter removed).
func TextColWidth(termW int) int { return ContentWidth(termW) - GutterVisualWidth }

// TextHeight is the number of visible text rows inside the main panel.
// Accounts for border (top+bottom = 2), title line (1), and status line (1).
func TextHeight(termH int) int {
	h := termH - 4
	if h < 1 {
		return 1
	}
	return h
}

// WordLayout records which wrapped line each word falls on.
type WordLayout struct {
	LineOf       []int // LineOf[wordIdx] = wrapped-line index (0-based)
	LineStart    []int // LineStart[lineIdx] = first word index on that line
	SourceLineOf []int // SourceLineOf[lineIdx] = 1-based source line number (shared across soft-wrapped continuations)
}

// Newline is the sentinel value inserted into the word list by main to
// represent a hard line break from the source file.
const Newline = "\n"

// ComputeLayout wraps words into lines of at most colW visible characters
// and returns the resulting WordLayout. Word entries equal to Newline are
// treated as forced line breaks; the sentinel occupies its own line slot so
// that word indices remain stable.
func ComputeLayout(words []string, colW int) WordLayout {
	if colW < 1 {
		colW = 1
	}
	lineOf := make([]int, len(words))
	lineStart := []int{0}
	li, lineW := 0, 0
	firstOnLine := true
	srcLine := 1
	sourceLineOf := []int{srcLine}
	for i, w := range words {
		if w == Newline {
			// Assign the sentinel to the current line (it won't be rendered).
			// Only the next word begins a new line.
			lineOf[i] = li
			li++
			srcLine++
			if i+1 < len(words) {
				lineStart = append(lineStart, i+1)
				sourceLineOf = append(sourceLineOf, srcLine)
			}
			lineW = 0
			firstOnLine = true
			continue
		}
		wl := len(w) // Latin text is ASCII
		if firstOnLine {
			lineW = wl
			firstOnLine = false
		} else if lineW+1+wl > colW {
			li++
			lineStart = append(lineStart, i)
			sourceLineOf = append(sourceLineOf, srcLine) // soft wrap: same source line
			lineW = wl
		} else {
			lineW += 1 + wl
		}
		lineOf[i] = li
	}
	return WordLayout{LineOf: lineOf, LineStart: lineStart, SourceLineOf: sourceLineOf}
}
