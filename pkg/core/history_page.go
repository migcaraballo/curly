package core

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HistoryPage struct {
	historyItems *[]CurlHistoryItem

	/* ui components */
	historyList    *tview.List
	Layout         *tview.Grid
	focusHandler   func(p tview.Primitive)
	defocusHandler func()
}

func NewHistoryPage(items *[]CurlHistoryItem) *HistoryPage {
	h := &HistoryPage{
		historyItems: items,
		historyList:  tview.NewList(),
		Layout:       tview.NewGrid(),
	}
	h.initUI()
	return h
}

func (hp *HistoryPage) initUI() error {
	//hp.Layout.SetTitle("| Curl History |")
	//hp.Layout.SetTitleAlign(tview.AlignLeft)
	//hp.Layout.SetTitleColor(tcell.ColorYellow)
	hp.Layout.SetBorders(true)
	hp.Layout.SetBordersColor(tcell.ColorDodgerBlue)
	hp.Layout.SetRows(0)
	hp.Layout.SetColumns(0)

	for i, v := range *hp.historyItems {
		hp.historyList.AddItem(
			fmt.Sprintf("%s - %s", v.GetFormattedDate(), v.Request.Url),
			fmt.Sprintf("%s - %s", v.Request.Method, v.Request.TlsVer),
			rune(i+1), nil)
	}

	hp.historyList.SetTitle("| Curl History |")
	hp.historyList.SetTitleColor(tcell.ColorYellow)
	hp.historyList.SetTitleAlign(tview.AlignLeft)
	hp.historyList.SetBorder(true)
	hp.historyList.SetBorderColor(tcell.ColorDodgerBlue)

	hp.historyList.SetHighlightFullLine(true)
	hp.historyList.SetSelectedTextColor(tcell.ColorBlack)
	hp.historyList.SetSelectedBackgroundColor(tcell.ColorYellow)

	hp.Layout.AddItem(hp.historyList, 0, 0, 1, 1, 0, 0, true)

	// setup esc capture
	hp.Layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			hp.defocusHandler()
		}

		return event
	})

	return nil
}

func (h *HistoryPage) SetFocusHandler(f func(p tview.Primitive)) {
	h.focusHandler = f
}

func (h HistoryPage) GetLoadFocusItem() tview.Primitive {
	return h.historyList
}

func (h *HistoryPage) SetDefocusHandler(f func()) {
	h.defocusHandler = f
}
