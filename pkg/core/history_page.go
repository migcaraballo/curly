package core

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HistoryPage struct {
	historyItems *[]CurlHistoryItem

	/* ui components */
	pageStage        *tview.Flex
	historyList      *tview.List
	Layout           *tview.Grid
	resultDetailView *tview.TextView
	focusHandler     func(p tview.Primitive)
	defocusHandler   func()
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
	hp.Layout.SetBorders(true)
	hp.Layout.SetBordersColor(tcell.ColorDodgerBlue)
	hp.Layout.SetRows(0, 1)
	hp.Layout.SetColumns(0)

	for _, v := range *hp.historyItems {
		hp.historyList.AddItem(hp.FormatLineItem(&v), "", 0, func() {
			tl := *hp.historyItems
			hi := tl[hp.historyList.GetCurrentItem()]
			hp.showHistoryItemDetails(&hi)
		})
	}

	hp.historyList.SetTitle("| Curl History |")
	hp.historyList.SetTitleColor(tcell.ColorYellow)
	hp.historyList.SetTitleAlign(tview.AlignLeft)
	hp.historyList.SetBorder(true)
	hp.historyList.SetBorderColor(tcell.ColorDodgerBlue)
	hp.historyList.SetHighlightFullLine(true)
	hp.historyList.SetSelectedTextColor(tcell.ColorBlack)
	hp.historyList.SetSelectedBackgroundColor(tcell.ColorYellow)

	/* page stage */
	hp.pageStage = tview.NewFlex()
	hp.pageStage.AddItem(hp.historyList, 0, 1, false)

	hp.Layout.AddItem(hp.pageStage, 0, 0, 1, 1, 0, 0, true)
	hp.Layout.AddItem(tview.NewTextView().SetText(NAV_TEXT).SetTextAlign(tview.AlignLeft), 1, 0, 1, 1, 0, 0, false)

	// setup esc capture
	hp.Layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			hp.defocusHandler()
		}

		return event
	})

	return nil
}

func (h *HistoryPage) FormatLineItem(c *CurlHistoryItem) string {
	return fmt.Sprintf("%s [m: %s] [t: %s] %s", c.GetFormattedDate(), c.Request.Method, c.Request.TlsVer, c.Request.Url)
}

func (h *HistoryPage) showHistoryItemDetails(hi *CurlHistoryItem) {
	d := tview.NewTextView()
	//d.SetBorder(true)
	d.SetDynamicColors(true)
	d.SetRegions(true)
	d.SetText(*hi.CurlResult)

	d.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlZ {
			h.switchStage(h.historyList)
			return nil
		}

		return event
	})

	h.switchStage(d)
}

func (h *HistoryPage) switchStage(p tview.Primitive) {
	h.pageStage.Clear()
	h.pageStage.AddItem(p, 0, 1, false)
	h.focusHandler(p)
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
