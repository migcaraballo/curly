package core

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

type App struct {
	currPage     *tview.Primitive
	curlyService *CurlyService

	/* UI components */
	tApp        *tview.Application
	mainLayout  *tview.Grid
	stage       *tview.Flex
	menu        *tview.List
	welcomePage *tview.TextView
	curlyPage   *CurlyPage
	HistoryPage *HistoryPage
}

func NewApp() (*App, error) {
	a := &App{}
	a.curlyService = NewCurlyService()

	//check for curl
	fmt.Println("Checking Curl")
	curlPath, err := a.curlyService.CheckCurl()
	if err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		// todo: need some curl interface for use with factory to toggle between real curl and golang curl
		fmt.Println("Found curl:", curlPath)
	}

	a.tApp = tview.NewApplication()

	// create the main UI
	fmt.Println("Initializing Curly UI...")
	a.mainLayout = a.createNewLayout()

	if err := a.tApp.SetRoot(a.mainLayout, true).SetFocus(a.mainLayout).Run(); err != nil {
		panic(err)
	}

	fmt.Println("exiting curly, thanks for using me.")

	return a, nil
}

func (a *App) createNewLayout() *tview.Grid {
	lo := tview.NewGrid()
	lo.SetBorders(true)
	lo.SetBackgroundColor(tcell.ColorDeepSkyBlue)
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
	a.welcomePage = tview.NewTextView().SetText(WelcomePageText).SetTextAlign(tview.AlignLeft).SetDynamicColors(true).SetRegions(true)
	a.welcomePage.SetBorder(true)
	a.welcomePage.SetBackgroundColor(tcell.ColorDarkSlateGray)
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
			a.curlyPage = NewCurlyPage()
			a.curlyPage.SetFocusHandler(func(p tview.Primitive) {
				a.tApp.SetFocus(p)
			})

			a.curlyPage.SetDeFocusHandler(func() {
				a.tApp.SetFocus(a.menu)
			})

			a.curlyPage.SetCurlCallHandler(func(creq *CurlRequest) string {
				res, _ := a.curlyService.ExecuteCurlCall(creq)
				return res
			})
		}

		a.setStage(a.curlyPage.GetMainPage(), false)
		a.tApp.SetFocus(a.curlyPage.GetItemForFocus())
	})

	m.AddItem("History", "", 'g', func() {
		a.HistoryPage = NewHistoryPage(a.curlyService.historyQueue)
		a.HistoryPage.SetFocusHandler(func(p tview.Primitive) {
			a.tApp.SetFocus(p)
		})
		a.HistoryPage.SetDefocusHandler(func() {
			a.tApp.SetFocus(m)
		})

		a.setStage(a.HistoryPage.Layout, true)
		a.tApp.SetFocus(a.HistoryPage.GetLoadFocusItem())
	})

	m.AddItem("Quit", "", 'q', func() {
		ed := tview.NewModal()
		ed.SetTextColor(tcell.ColorBlack)
		ed.SetText("Sure you wanna Quit?")
		ed.AddButtons([]string{"Yes", "No"})
		ed.SetDoneFunc(func(bidx int, blbl string) {
			if strings.EqualFold("yes", blbl) {
				a.tApp.Stop()
			} else {
				a.tApp.SetRoot(a.mainLayout, true)
				a.tApp.SetFocus(m)
			}
		})

		a.tApp.SetRoot(ed, false).SetFocus(ed)
	})

	return m
}
