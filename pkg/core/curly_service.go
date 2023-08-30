package core

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

var (
	CURL_PATHS = []string{
		"/bin/curl",
		"/usr/bin/curl",
		"/usr/local/bin/curl",
	}
)

type CurlHistoryItem struct {
	Request    *CurlRequest
	CurlResult *string
}

type CurlyService struct {
	historyQueue    *[]CurlHistoryItem
	useableCurlPath string
}

func NewCurlyService() *CurlyService {
	return &CurlyService{
		historyQueue: &[]CurlHistoryItem{},
	}
}

func (cs *CurlyService) AddResult(creq *CurlRequest, result *string) {
	*cs.historyQueue = append(*cs.historyQueue, CurlHistoryItem{creq, result})
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

func (cs CurlyService) ExecuteCurlCall(creq *CurlRequest) (string, error) {
	return cs.curlIt(creq)
}

func (cs *CurlyService) CheckCurl() (string, error) {
	for _, c := range CURL_PATHS {
		if cpath, err := exec.LookPath(c); err == nil {
			cs.useableCurlPath = cpath
			return cpath, nil
		}
	}

	return "", errors.New("CURL was not found on this computer")
}

// https://api.dictionaryapi.dev/api/v2/entries/en/test
func (cs CurlyService) curlIt(cr *CurlRequest) (string, error) {
	vArgs := []string{}
	vArgs = append(vArgs, "-v")
	vArgs = append(vArgs, cr.Url)

	if len(cr.Headers) > 0 {
		vArgs = append(vArgs, cr.GetHeadersString())
	}
	if len(cr.QsParams) > 0 {
		vArgs = append(vArgs, cr.GetQsParamString())
	}
	if len(cr.Body) > 0 {
		vArgs = append(vArgs, cr.GetBodyString())
	}

	cmd := exec.Command(cs.useableCurlPath, vArgs...)

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	res := string(stdout)
	cs.AddResult(cr, &res)

	return res, nil
}