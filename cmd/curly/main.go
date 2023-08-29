package main

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"os/exec"
	"strings"
)

var (
	CURL_PATHS = []string{
		"/bin/curl",
		"/usr/bin/curl",
		"/usr/local/bin/curl",
	}
	useableCurlPath string
	currPage        *tview.Primitive

	/* UI components */
	app        *tview.Application
	mainLayout *tview.Grid
	stage      *tview.Flex
	//pages       *tview.Pages
	menu        = tview.NewList()
	welcomePage *tview.TextView
)

func main() {
	// check for curl
	//fmt.Println("Checking Curl")
	//curlPath, err := checkCurl()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//} else {
	//	// future: need some curl interface for use with factory to toggle between real curl and golang curl
	//	fmt.Println("Found curl:", curlPath)
	//	useableCurlPath = curlPath
	//}

	app = tview.NewApplication()
	//pages = tview.NewPages()

	// create the main UI
	fmt.Println("Initializing Curly UI...")
	mainPage := createNewLayout()
	//pages.AddPage("main", mainPage, false, true)

	//if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
	if err := app.SetRoot(mainPage, true).SetFocus(mainPage).Run(); err != nil {
		panic(err)
	}

	fmt.Println("exiting curly, thanks for using me.")
}

func createNewLayout() *tview.Grid {
	lo := tview.NewGrid()
	lo.SetBorders(true)
	//lo.SetBackgroundColor(tcell.ColorBlack)
	//lo.SetBordersColor(tcell.ColorDodgerBlue)
	//lo.SetBorderColor(tcell.ColorDodgerBlue)

	// 1 row
	lo.SetRows(0)

	// left nav | right content
	lo.SetColumns(26, 0)

	// welcome page
	welcomePage = tview.NewTextView().SetText("Welcome To Curly!").SetTextAlign(tview.AlignLeft)
	welcomePage.SetBorder(true)
	welcomePage.SetTextColor(tcell.ColorWhite)
	welcomePage.SetBorderColor(tcell.ColorDarkCyan)
	lo.AddItem(welcomePage, 0, 1, 1, 1, 0, 0, false)

	// add menu
	menu = createPageMenu()
	lo.AddItem(menu, 0, 0, 1, 1, 0, 0, true)

	return lo
}

func createPageMenu() *tview.List {
	m := tview.NewList()
	m.ShowSecondaryText(false)
	m.SetBorder(true)
	m.SetMainTextColor(tcell.ColorWhite)

	m.AddItem("Curl It!", "", 'c', func() {
		//stage.Clear()
		//stage.AddItem(welcomeText, 0, 1, true)
	})

	m.AddItem("History", "", 'g', func() {
		//stage.Clear()
		//initApiGetScreen()
		//stage.AddItem(apiGetScreen, 0, 1, true)
	})

	m.AddItem("Quit", "", 'q', func() {
		ed := tview.NewModal()
		ed.SetTextColor(tcell.ColorBlack)
		ed.SetText("Sure you wanna Quit?")
		ed.AddButtons([]string{"Yes", "No"})
		ed.SetDoneFunc(func(bidx int, blbl string) {
			if strings.EqualFold("yes", blbl) {
				app.Stop()
			}

			// how to get back to previous page? use 'currPage' VAR

			//else {
			//	app.SetRoot(grd, true)
			//	app.SetFocus(m)
			//}
		})

		app.SetRoot(ed, false).SetFocus(ed)
	})

	return m
}

func checkCurl() (string, error) {
	for _, c := range CURL_PATHS {
		if cpath, err := exec.LookPath(c); err == nil {
			return cpath, nil
		}
	}

	return "", errors.New("CURL was not found on this computer")
}

func curlIt(args ...string) string {
	cmd := exec.Command("curl", "-v", "https://api.dictionaryapi.dev/api/v2/entries/en/test")

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("ERR:", err)
		return ""
	}

	return string(stdout)
}
