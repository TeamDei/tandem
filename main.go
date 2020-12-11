package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	COLUMNS int // Columns
	FILE string // Path to file to open
)

func init() {
	flag.StringVar(&FILE, "file", "", "path to file")
	flag.IntVar(&COLUMNS, "cols", 10, "columns to display")

	flag.Parse()
}

func main() {
	app := tview.NewApplication()

	// Create a table object
	table := tview.NewTable()
	// Create a pages object
	pages := tview.NewPages()

	// Load data from a file
	d, err := LoadFile(FILE)
	if err != nil {
		Handle(err)
	}

	// Replace line endings with spaces
	d = strings.Replace(strings.Replace(d, "\r", "", -1), "\n", " ", -1)
	// Populate table with data from the selected file
	data := strings.Split(tview.Escape(d), " ")
	r, c := 0, 0
	for i := 0; i<len(data); i++ {
		if c == 0 {
			table.SetCell(r, 0, tview.NewTableCell(fmt.Sprintf("%d", r+1)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
			table.SetCell(r, 1, tview.NewTableCell(data[i]).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
			c++
		} else {
			table.SetCell(r, c+1,
				tview.NewTableCell(data[i]).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
			if c == COLUMNS {
				c = 0; r++
			} else { c++ }
		}
	}

	// Add event handler to handle the selection of a table element
	table.SetSelectedFunc(func(row int, column int) {
		// Create a new flex wrapper
		flex := tview.NewFlex()
		// Create a new list object
		definition := tview.NewList()
		// Query the API
		r, err := GenAPI(table.GetCell(row, column).Text)
		if err != nil {
			panic(err)
		}
		// Add analyses to the list
		for _, a := range r.Analyses {
			definition.AddItem(a.String(), "", '~', nil)
		}
		// Add exit option to dialog
		definition.AddItem("Exit", "Close dialog", 'e', func() {
			pages.RemovePage("define")
			pages.SwitchToPage("main")
		})
		// Add the list to the flex
		flex.AddItem(nil, 0, 1, false).
			AddItem(definition, 0, 2, true).
			AddItem(nil, 0, 1, false).
			SetBorder(true).
			SetTitle("Tandem - Analyses")

		pages.AddPage("define", flex, true, false)
		pages.SwitchToPage("define")
	})
	// Add event handler to handle a keypress
	table.Select(0, 0).SetDoneFunc(func(key tcell.Key) {
		// If ESC is pressed, open the exit dialog
		if key == tcell.KeyEscape {
			pages.SwitchToPage("exitdialog")
		}
	})

	// Add a border and title to the table
	table.SetSelectable(true, true).SetBorder(true).SetTitle("Tandem")
	// Add the table to the pages object
	pages.AddPage("main", table, true, true)
	// Add an exit dialog to the pages object
	pages.AddPage("exitdialog", tview.NewModal().
		SetText("Do you want to quit tandem?").AddButtons([]string{"Cancel", "Quit"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				app.Stop()
			} else {
				pages.SwitchToPage("main")
			}
		}), true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		Handle(err)
	}
}

// Handle an error
func Handle(e error) {
	fmt.Println("error,", e)
	os.Exit(1)
}

// LoadFile returns a string containing the contents of the file
func LoadFile(path string) (string, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(d), nil
}
