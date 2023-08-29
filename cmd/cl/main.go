package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	app         = tview.NewApplication()
	grd         = tview.NewGrid()
	navList     = tview.NewList()
	stage       = tview.NewFlex()
	welcomeText = tview.NewTextView()

	/* screens */
	apiGetScreen = tview.NewGrid()
)

func main() {
	initLayout()

	if err := app.SetRoot(grd, true).SetFocus(grd).Run(); err != nil {
		panic(err)
	}
}

func initLayout() {
	// main layout
	grd.SetRows(0)
	grd.SetColumns(26, 0)
	grd.SetBorders(true)
	grd.SetBackgroundColor(tcell.ColorBlack)

	// menu
	initMenu()

	welcomeText.SetText("The stage is set.").SetTextAlign(tview.AlignCenter)
	stage.SetBackgroundColor(tcell.ColorDarkViolet)
	stage.SetBorder(true)

	welcomeText.SetBackgroundColor(tcell.ColorBlack)
	stage.AddItem(welcomeText, 0, 1, false)

	grd.AddItem(navList, 0, 0, 1, 1, 0, 0, true)
	grd.AddItem(stage, 0, 1, 1, 1, 0, 0, false)
}

func initApiGetScreen() {
	apiGetScreen = tview.NewGrid()
	apiGetScreen.SetRows(10, 0)
	apiGetScreen.SetColumns(0)
	apiGetScreen.SetBorders(true)

	// results conatainer
	results := tview.NewTextView().SetText("Results").SetTextAlign(tview.AlignCenter)

	// setup the text box for the call
	f := tview.NewForm()
	f.SetBorder(true)
	f.SetTitle("Query the Open Dictionary API")
	f.SetTitleAlign(tview.AlignCenter)
	f.SetFieldBackgroundColor(tcell.ColorLightGray)
	f.SetFieldTextColor(tcell.ColorWhite)
	f.SetButtonTextColor(tcell.ColorBlack)
	f.SetBorderColor(tcell.ColorGreen)
	f.SetBorder(true)
	f.SetButtonBackgroundColor(tcell.ColorGreen)

	f.AddTextView("Note:", "Enter a word to lookup", 0, 2, true, false)

	f.AddTextArea("Term:", "", 0, 1, 0, nil)

	f.AddButton("Submit", func() {
		tf := f.GetFormItemByLabel("Term:").(*tview.TextArea)
		q := tf.GetText()
		q = strings.TrimSpace(q)
		results.SetTextAlign(tview.AlignLeft)
		q = url.QueryEscape(q)
		ress := handleSearch(q)
		results.SetText(
			fmt.Sprintf("Query: %s\n\nTotal Results: %d\n=======================================================\n\n%s",
				q,
				len(ress),
				ress))
		app.SetFocus(results)

		results.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEsc {
				app.SetFocus(f)
			}

			return event
		})

		results.ScrollTo(0, 0)
	})

	f.AddButton("Reset", func() {
		tf := f.GetFormItemByLabel("Term:").(*tview.TextArea)
		tf.SetText("", false)
		results.SetText("")
		app.SetFocus(f.GetFormItemByLabel("Term:"))
		return
	})

	f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			stage.Clear()
			app.SetFocus(navList)
			stage.AddItem(welcomeText, 0, 1, true)
			return nil
		}

		return event
	})

	apiGetScreen.AddItem(f, 0, 0, 1, 1, 0, 0, true)
	apiGetScreen.AddItem(results, 1, 0, 1, 1, 0, 0, false)

	app.SetFocus(f)
}

func initMenu() {
	navList.ShowSecondaryText(false)
	navList.SetBorder(true)

	navList.AddItem("Home", "", 'h', func() {
		stage.Clear()
		stage.AddItem(welcomeText, 0, 1, true)
	})

	navList.AddItem("Urban Dictionary", "", 'g', func() {
		stage.Clear()
		initApiGetScreen()
		stage.AddItem(apiGetScreen, 0, 1, true)
	})

	navList.AddItem("Quit", "", 'q', func() {
		ed := tview.NewModal()
		ed.SetTextColor(tcell.ColorBlack)
		ed.SetText("Sure you wanna Quit?")
		ed.AddButtons([]string{"Yes", "No"})
		ed.SetDoneFunc(func(bidx int, blbl string) {
			if strings.EqualFold("yes", blbl) {
				app.Stop()
			} else {
				app.SetRoot(grd, true)
				app.SetFocus(navList)
			}
		})

		app.SetRoot(ed, false).SetFocus(ed)
	})
}

func handleSearch(q string) string {
	//url := fmt.Sprintf("https://api.urbandictionary.com/v0/define?term=%s", q)
	url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", q)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return ""
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}

	defer res.Body.Close()

	//var tms Terms
	//err = json.NewDecoder(res.Body).Decode(&tms)
	var dres DictResp
	err = json.NewDecoder(res.Body).Decode(&dres)
	if err != nil {
		panic(err)
	}

	s := ""
	for _, d := range dres {
		for _, m := range d.Meanings {
			s += fmt.Sprintf("%v\n\n", m.Definitions[0])
			s += "-------------------------------------------------------\n"
		}
	}
	return s
}

type Terms struct {
	List []struct {
		Definition  string    `json:"definition"`
		Permalink   string    `json:"permalink"`
		ThumbsUp    int       `json:"thumbs_up"`
		Author      string    `json:"author"`
		Word        string    `json:"word"`
		Defid       int       `json:"defid"`
		CurrentVote string    `json:"current_vote"`
		WrittenOn   time.Time `json:"written_on"`
		Example     string    `json:"example"`
		ThumbsDown  int       `json:"thumbs_down"`
	} `json:"list"`
}

type DictResp []DictRespElement

type DictRespElement struct {
	Word       string     `json:"word"`
	Phonetics  []Phonetic `json:"phonetics"`
	Meanings   []Meaning  `json:"meanings"`
	License    License    `json:"license"`
	SourceUrls []string   `json:"sourceUrls"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
	Synonyms     []string     `json:"synonyms"`
	Antonyms     []string     `json:"antonyms"`
}

type Definition struct {
	Definition string        `json:"definition"`
	Synonyms   []interface{} `json:"synonyms"`
	Antonyms   []interface{} `json:"antonyms"`
	Example    *string       `json:"example,omitempty"`
}

type Phonetic struct {
	Audio     string  `json:"audio"`
	SourceURL string  `json:"sourceUrl"`
	License   License `json:"license"`
	Text      *string `json:"text,omitempty"`
}
