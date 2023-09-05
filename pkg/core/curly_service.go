package core

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var (
	CURL_PATHS = []string{
		"/bin/curl",
		"/usr/bin/curl",
		"/usr/local/bin/curl",
	}
)

type CurlyService struct {
	historyQueue    *[]CurlHistoryItem
	useableCurlPath string
	curlClient      Curler
}

func NewCurlyService() *CurlyService {
	cs := &CurlyService{
		historyQueue: &[]CurlHistoryItem{},
	}

	//check for curl
	fmt.Println("Checking Curl")
	curlPath, err := cs.checkCurl()

	if err != nil {
		fmt.Printf("Curl was not found, using Curly client: %v\n", err)
		cs.SetCurlClient(false)
	} else {
		fmt.Println("Found curl:", curlPath)
		cs.useableCurlPath = curlPath
		cs.SetCurlClient(true)
	}

	return cs
}

func (cs *CurlyService) AddResult(creq *CurlRequest, result *string) {
	// prepend behavior
	*cs.historyQueue = append([]CurlHistoryItem{CurlHistoryItem{creq, result}}, *cs.historyQueue...)
}

func (cs *CurlyService) GetCurlHistory() *[]CurlHistoryItem {
	return cs.historyQueue
}

func (cs CurlyService) GetCurlHistoryString() string {
	b := bytes.Buffer{}
	for i, ch := range *cs.historyQueue {
		b.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, ch.Request.RequestDate.Format(time.RFC822), ch.Request.Url))
	}

	return b.String()
}

func (cs *CurlyService) ExecuteCurlCall(creq *CurlRequest) (string, error) {
	res, err := cs.curlClient.CurlCall(creq)
	res = strings.TrimSpace(res)
	res += "\n\n<<< EOF >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
	cs.AddResult(creq, &res)
	return res, err
}

func (cs *CurlyService) SetCurlClient(native bool) {
	if native {
		cs.curlClient = NewNativeCurlClient(cs.useableCurlPath)
	} else {
		cs.curlClient = NewGoCurlClient()
	}
}

func (cs *CurlyService) checkCurl() (string, error) {
	for _, c := range CURL_PATHS {
		if cpath, err := exec.LookPath(c); err == nil {
			return cpath, nil
		}
	}

	return "", errors.New("CURL was not found on this computer")
}
