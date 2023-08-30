package core

import (
	"curly/internal/page"
	"curly/internal/service"
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
)

type App struct {
	useableCurlPath string
	currPage        *tview.Primitive
	curlyService    *service.CurlyService

	/* UI components */
	tApp        *tview.Application
	mainLayout  *tview.Grid
	stage       *tview.Flex
	menu        *tview.List
	welcomePage *tview.TextView
	curlyPage   *page.CurlyPage
}

func NewApp() (*App, error) {
	a := &App{}
	a.curlyService = service.NewCurlyService()

	//check for curl
	fmt.Println("Checking Curl")
	curlPath, err := a.checkCurl()
	if err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		// future: need some curl interface for use with factory to toggle between real curl and golang curl
		fmt.Println("Found curl:", curlPath)
		a.useableCurlPath = curlPath
	}

	a.tApp = tview.NewApplication()

	// create the main UI
	fmt.Println("Initializing Curly UI...")
	mainPage := a.createNewLayout()

	if err := a.tApp.SetRoot(mainPage, true).SetFocus(mainPage).Run(); err != nil {
		panic(err)
	}

	fmt.Println("exiting curly, thanks for using me.")

	return a, nil
}

func (a *App) createNewLayout() *tview.Grid {
	lo := tview.NewGrid()
	lo.SetBorders(true)
	lo.SetBackgroundColor(tcell.ColorBlack)
	lo.SetBordersColor(tcell.ColorDodgerBlue)
	lo.SetBorderColor(tcell.ColorDodgerBlue)

	// 1 row
	lo.SetRows(0)

	// left nav | right content
	lo.SetColumns(26, 0)

	// add menu
	a.menu = a.createPageMenu()
	lo.AddItem(a.menu, 0, 0, 1, 1, 0, 0, true)

	// welcome page
	a.welcomePage = tview.NewTextView().SetText("Welcome To Curly!").SetTextAlign(tview.AlignLeft)
	a.welcomePage.SetBorder(true)
	a.welcomePage.SetTextColor(tcell.ColorWhite)
	a.welcomePage.SetBorderColor(tcell.ColorDodgerBlue)

	a.stage = tview.NewFlex()
	a.setStage(a.welcomePage, false)
	lo.AddItem(a.stage, 0, 1, 1, 1, 0, 0, false)

	return lo
}

func (a *App) setStage(tp tview.Primitive, focus bool) {
	a.stage.Clear()
	a.stage.AddItem(tp, 0, 1, focus)
}

func (a *App) createPageMenu() *tview.List {
	m := tview.NewList()
	m.ShowSecondaryText(false)
	m.SetBorder(true)
	m.SetMainTextColor(tcell.ColorWhite)

	m.AddItem("Curl It!", "", 'c', func() {
		if a.curlyPage == nil {
			a.curlyPage = page.NewCurlyPage()
			a.curlyPage.SetFocusHandler(func(p tview.Primitive) {
				a.tApp.SetFocus(p)
			})

			a.curlyPage.SetDeFocusHandler(func() {
				a.tApp.SetFocus(a.menu)
			})

			a.curlyPage.SetCurlCallHandler(func(creq *page.CurlRequest, result *string) {
				a.curlyService.AddResult(creq, result)
			})
		}

		a.setStage(a.curlyPage.GetMainPage(), false)
		a.tApp.SetFocus(a.curlyPage.GetItemForFocus())
	})

	m.AddItem("History", "", 'g', func() {
		tmpHistPage := tview.NewTextView().SetText("Temp History Page")
		tmpHistPage.SetBorder(true)
		tmpHistPage.SetBackgroundColor(tcell.ColorBlack)
		tmpHistPage.SetBorderColor(tcell.ColorDodgerBlue)

		tmpHistPage.SetText(a.curlyService.GetCurlHistoryString())

		a.setStage(tmpHistPage, false)
	})

	m.AddItem("Quit", "", 'q', func() {
		ed := tview.NewModal()
		ed.SetTextColor(tcell.ColorBlack)
		ed.SetText("Sure you wanna Quit?")
		ed.AddButtons([]string{"Yes", "No"})
		ed.SetDoneFunc(func(bidx int, blbl string) {
			if strings.EqualFold("yes", blbl) {
				a.tApp.Stop()
			}

			// how to get back to previous page? use 'currPage' VAR

			//else {
			//	app.SetRoot(grd, true)
			//	app.SetFocus(m)
			//}
		})

		a.tApp.SetRoot(ed, false).SetFocus(ed)
	})

	return m
}

func (a *App) checkCurl() (string, error) {
	for _, c := range CURL_PATHS {
		if cpath, err := exec.LookPath(c); err == nil {
			return cpath, nil
		}
	}

	return "", errors.New("CURL was not found on this computer")
}

func (a *App) curlIt(args ...string) string {
	cmd := exec.Command("curl", "-v", "https://api.dictionaryapi.dev/api/v2/entries/en/test")

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("ERR:", err)
		return ""
	}

	return string(stdout)
}
