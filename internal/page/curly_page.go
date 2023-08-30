package page

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type CurlRequest struct {
	RequestDate time.Time
	Url         string
	tlsVer      string
	method      string
	headers     string
	qsParams    string
	body        string
}

type CurlyPage struct {
	mainGrid        *tview.Grid
	form            *tview.Form
	results         *tview.TextView
	defocusHandler  func()
	focusHandler    func(t tview.Primitive)
	curlCallHandler func(curlRequest *CurlRequest, result *string)
}

func NewCurlyPage() *CurlyPage {
	cp := &CurlyPage{
		form:     tview.NewForm(),
		mainGrid: tview.NewGrid(),
		results:  tview.NewTextView(),
	}

	cp.initUI()
	return cp
}

func (cp CurlyPage) GetMainPage() tview.Primitive {
	return cp.mainGrid
}

func (cp *CurlyPage) initUI() {
	cp.form.SetBorderColor(tcell.ColorDodgerBlue)
	cp.form.SetItemPadding(1)
	cp.form.SetBorder(true)
	cp.form.SetTitle("| Call a URL |")
	cp.form.SetTitleAlign(tview.AlignLeft)

	cp.form.SetFieldBackgroundColor(tcell.ColorWhite)
	cp.form.SetFieldTextColor(tcell.ColorBlack)
	cp.form.SetLabelColor(tcell.ColorWhite)

	cp.form.AddInputField("URL", "", 0, nil, nil)
	cp.form.AddDropDown("Method", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, 0, nil)
	cp.form.AddDropDown("TLS", []string{"1.0", "1.1", "1.2", "1.3"}, 0, nil)
	cp.form.AddTextArea("Query Params", "", 0, 3, 0, nil)
	cp.form.AddTextArea("Headers", "", 0, 3, 0, nil)
	cp.form.AddTextArea("Body", "", 0, 5, 0, nil)

	// buttons
	cp.form.SetButtonBackgroundColor(tcell.ColorDodgerBlue)
	cp.form.SetButtonTextColor(tcell.ColorBlack)

	// Call Button + behavior
	cp.form.AddButton("Call", func() {
		creq := cp.GetCurlRequst()
		cp.results.Clear().SetText(fmt.Sprintf("%+v\n", creq))
		cp.focusHandler(cp.results)

		res := cp.results.GetText(false)
		cp.curlCallHandler(creq, &res)
		return
	})

	// clear button + behavior
	cp.form.AddButton("Clear", nil)

	cp.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			cp.defocusHandler()
		}

		return event
	})

	// results view
	cp.results.SetBorder(true)
	cp.results.SetTitle("| Results |")
	cp.results.SetTitleAlign(tview.AlignLeft)
	cp.results.SetBorderColor(tcell.ColorDodgerBlue)
	cp.results.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			cp.focusHandler(cp.form)
		}

		return event
	})

	// setup grid layout
	cp.mainGrid.SetRows(25, 0)
	cp.mainGrid.SetColumns(0)
	cp.mainGrid.AddItem(cp.form, 0, 0, 1, 1, 0, 0, false)
	cp.mainGrid.AddItem(cp.results, 1, 0, 1, 1, 0, 0, false)
}

func (cp CurlyPage) GetCurlRequst() *CurlRequest {
	_, tVal := cp.form.GetFormItemByLabel("TLS").(*tview.DropDown).GetCurrentOption()
	_, mVal := cp.form.GetFormItemByLabel("Method").(*tview.DropDown).GetCurrentOption()

	return &CurlRequest{
		RequestDate: time.Now(),
		Url:         cp.form.GetFormItemByLabel("URL").(*tview.InputField).GetText(),
		tlsVer:      tVal,
		method:      mVal,
		headers:     cp.form.GetFormItemByLabel("Headers").(*tview.TextArea).GetText(),
		qsParams:    cp.form.GetFormItemByLabel("Query Params").(*tview.TextArea).GetText(),
		body:        cp.form.GetFormItemByLabel("Body").(*tview.TextArea).GetText(),
	}
}

func (cp *CurlyPage) SetCurlCallHandler(f func(creq *CurlRequest, result *string)) {
	cp.curlCallHandler = f
}

func (cp *CurlyPage) SetFocusHandler(f func(p tview.Primitive)) {
	cp.focusHandler = f
}

func (cp *CurlyPage) SetDeFocusHandler(f func()) {
	cp.defocusHandler = f
}

func (cp CurlyPage) GetItemForFocus() tview.Primitive {
	return cp.form
}
