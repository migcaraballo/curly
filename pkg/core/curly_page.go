package core

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
	"time"
)

type CurlyPage struct {
	mainGrid        *tview.Grid
	form            *tview.Form
	results         *tview.TextView
	defocusHandler  func()
	focusHandler    func(t tview.Primitive)
	curlCallHandler func(curlRequest *CurlRequest) string
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
	cp.form.AddDropDown("TLS", []string{"1.0", "1.1", "1.2", "1.3", "1.4"}, 0, nil)
	cp.form.AddTextArea("Query Params", "", 0, 3, 0, nil)
	cp.form.AddTextArea("Headers", "", 0, 3, 0, nil)
	cp.form.AddTextArea("Body", "", 0, 5, 0, nil)

	cp.ResetForm()

	// buttons
	cp.form.SetButtonBackgroundColor(tcell.ColorDodgerBlue)
	cp.form.SetButtonTextColor(tcell.ColorBlack)

	// Call Button + behavior
	cp.form.AddButton("Call", func() {
		cp.results.Clear() // clear results before others pop-in
		creq := cp.GetCurlRequst()
		res := cp.curlCallHandler(creq)

		cp.results.Clear().SetText(fmt.Sprintf("%s\n", res))
		cp.results.SetScrollable(true)
		cp.results.ScrollToBeginning()

		cp.focusHandler(cp.results)
		return
	})

	// clear button + behavior
	cp.form.AddButton("Clear", func() {
		cp.ResetForm()
	})

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

	// setup grid Layout
	cp.mainGrid.SetRows(25, 0)
	cp.mainGrid.SetColumns(0)
	cp.mainGrid.AddItem(cp.form, 0, 0, 1, 1, 0, 0, false)
	cp.mainGrid.AddItem(cp.results, 1, 0, 1, 1, 0, 0, false)
}

func (cp CurlyPage) GetCurlRequst() *CurlRequest {
	_, tVal := cp.form.GetFormItemByLabel("TLS").(*tview.DropDown).GetCurrentOption()
	_, mVal := cp.form.GetFormItemByLabel("Method").(*tview.DropDown).GetCurrentOption()
	hdrs := cp.ParseMultiLine(cp.form.GetFormItemByLabel("Headers").(*tview.TextArea).GetText())
	qsp := cp.ParseMultiLine(cp.form.GetFormItemByLabel("Query Params").(*tview.TextArea).GetText())

	cr := NewCurlRequest()
	cr.RequestDate = time.Now()
	cr.Url = cp.form.GetFormItemByLabel("URL").(*tview.InputField).GetText()
	cr.TlsVer = tVal
	cr.Method = mVal
	cr.Headers = hdrs
	cr.QsParams = qsp
	cr.Body = cp.form.GetFormItemByLabel("Body").(*tview.TextArea).GetText()

	return cr
}

func (cp CurlyPage) ParseMultiLine(ml string) []string {
	ml = strings.TrimSpace(ml)

	sarr := strings.Split(ml, "\n")
	if len(sarr[0]) == 0 {
		return nil
	}

	return sarr
}

func (cp *CurlyPage) ResetForm() {
	cp.form.GetFormItemByLabel("TLS").(*tview.DropDown).SetCurrentOption(0)
	cp.form.GetFormItemByLabel("Method").(*tview.DropDown).SetCurrentOption(0)
	cp.form.GetFormItemByLabel("Headers").(*tview.TextArea).SetText("", true).SetPlaceholder("key:value\nkey:value").
		SetPlaceholderStyle(tcell.StyleDefault)
	cp.form.GetFormItemByLabel("Query Params").(*tview.TextArea).SetText("", true).SetPlaceholder("key=value\nkey=value").
		SetPlaceholderStyle(tcell.StyleDefault)
	cp.form.GetFormItemByLabel("URL").(*tview.InputField).SetText("")
	cp.form.GetFormItemByLabel("Body").(*tview.TextArea).SetText("", true)
}

func (cp *CurlyPage) SetCurlCallHandler(f func(creq *CurlRequest) string) {
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
