package page

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CurlRequest struct {
	Url      string
	tlsVer   string
	method   string
	headers  map[string]string
	qsParams map[string]string
}

type CurlyPage struct {
	form           *tview.Form
	app            *tview.Application
	defocusHandler func()
	//defocusItem    tview.Primitive
}

func NewCurlyPage() *CurlyPage {
	cp := &CurlyPage{
		form: tview.NewForm(),
	}

	cp.initUI()
	return cp
}

func (cp CurlyPage) GetPrimitvie() tview.Primitive {
	return cp.form
}

func (cp *CurlyPage) initUI() {
	cp.form.SetBorderColor(tcell.ColorDodgerBlue)
	cp.form.SetItemPadding(1)
	cp.form.SetBorder(true)
	cp.form.SetTitle("Call a URL")
	cp.form.SetTitleAlign(tview.AlignLeft)

	cp.form.SetFieldBackgroundColor(tcell.ColorWhite)
	cp.form.SetFieldTextColor(tcell.ColorBlack)
	cp.form.SetLabelColor(tcell.ColorWhite)

	cp.form.AddInputField("URL", "", 0, nil, nil)
	cp.form.AddDropDown("Method", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, 0, nil)
	cp.form.AddDropDown("TLS", []string{"1.0", "1.1", "1.2", "1.3"}, 0, nil)
	cp.form.AddTextArea("QS Params", "", 0, 3, 0, nil)
	cp.form.AddTextArea("Headers", "", 0, 3, 0, nil)
	cp.form.AddTextArea("Body", "", 0, 5, 0, nil)

	// buttons
	cp.form.SetButtonBackgroundColor(tcell.ColorDodgerBlue)
	cp.form.SetButtonTextColor(tcell.ColorBlack)
	cp.form.AddButton("Call", nil)
	cp.form.AddButton("Clear", nil)

	cp.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			cp.defocusHandler()
		}

		return event
	})
}

func (cp CurlyPage) GetCurlRequst() *CurlRequest {
	_, tVal := cp.form.GetFormItemByLabel("TLS").(*tview.DropDown).GetCurrentOption()
	_, mVal := cp.form.GetFormItemByLabel("Method").(*tview.DropDown).GetCurrentOption()

	return &CurlRequest{
		Url:      cp.form.GetFormItemByLabel("URL").(*tview.InputField).GetText(),
		tlsVer:   tVal,
		method:   mVal,
		headers:  make(map[string]string),
		qsParams: make(map[string]string),
	}
}

func (cp *CurlyPage) SetDeFocusHandler(f func()) {
	cp.defocusHandler = f
}

func (cp CurlyPage) GetItemForFocus() tview.Primitive {
	//return cp.form.GetFormItem(0)
	return cp.form
}
