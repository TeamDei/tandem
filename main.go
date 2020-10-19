package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)
	lorem := strings.Split("Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.", " ")
	cols, rows := 10, 40
	word := 0
	for r := 0; r < rows; r++ {
		table.SetCell(r, 0,
			tview.NewTableCell(fmt.Sprintf("%d", r+1)).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter))
		for c := 1; c < cols; c++ {
			table.SetCell(r, c,
				tview.NewTableCell(lorem[word]).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignCenter))
			word = (word + 1) % len(lorem)
		}
	}
	table.Select(0, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	})
	table.SetSelectable(true, true).SetBorder(true).SetTitle("Tandem")

	definition := tview.NewList()
	pages := tview.NewPages()
	pages.AddPage("main", table.SetSelectedFunc(func(row int, column int) {
		
		pages.SwitchToPage("define")
	}), true, true)
	pages.AddPage("define", definition.AddItem("Quit", "Press to exit", 'q', func() {
		pages.SwitchToPage("main")
	}), true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

// Do tests (optionally test the API too)
func Tests(API bool) {
	Test("noun");Test("pronoun");Test("verb");Test("adjective");Test("adverb")
	if API {
		TestAPI("nautam");TestAPI("noster");TestAPI("hortari");TestAPI("acri");TestAPI("diligenter")
	}
	fmt.Println("")
}

// Parse an example and print analyses
func Test(id string) {
	fmt.Printf("\n===%s test===\n", id)
	r, err := GenFile(path.Join("test_data", "example_"+id+".xml"))
	if err != nil {
		panic(err)
	}
	for _, a := range r.Analyses {
		fmt.Println(a.String())
	}
}

// Parse example words from the actual API and print analyses
func TestAPI(word string) {
	fmt.Printf("\n===API TEST (%s)===\n", word)
	r, err := GenAPI(word)
	if err != nil {
		panic(err)
	}
	for _, a := range r.Analyses {
		fmt.Println(a.String())
	}
}
