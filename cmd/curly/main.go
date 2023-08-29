package main

import (
	"fmt"
	"github.com/rivo/tview"
)

func main() {
	fmt.Println("Initializing Curly UI...")
	app := tview.NewApplication()

	pages := tview.NewPages()

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}
