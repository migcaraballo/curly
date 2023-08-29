package main

import (
	"errors"
	"fmt"
	"github.com/rivo/tview"
	"log"
	"os/exec"
)

/*
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
*/
var (
	CURL_PATHS = []string{
		"/bin/curl",
		"/usr/bin/curl",
		"/usr/local/bin/curl",
	}
	useableCurlPath string

	/* UI components */
	mainLayout *tview.Grid
)

func main() {
	// check for curl
	fmt.Println("Checking Curl")
	curlPath, err := checkCurl()
	if err != nil {
		log.Fatalln(err)
		return
	} else {
		log.Println("Found curl:", curlPath)
		useableCurlPath = curlPath
	}

	// future: need some curl interface for use with factory to toggle between real curl and golang curl

	// create the main UI

	fmt.Println("EOF")
}

func createNewLayout() *tview.Grid {
	lo := tview.NewGrid()
	lo.SetBorder(true)
	lo.SetRows(0, 10)
	return lo
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
	//cmdRdr, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("ERR:", err)
		return ""
	}

	return string(stdout)
}
