package service

import (
	"bytes"
	"curly/internal/page"
	"fmt"
	"time"
)

type CurlHistoryItem struct {
	Request    *page.CurlRequest
	CurlResult *string
}

type CurlyService struct {
	historyQueue *[]CurlHistoryItem
}

func NewCurlyService() *CurlyService {
	return &CurlyService{
		historyQueue: &[]CurlHistoryItem{},
	}
}

func (cs *CurlyService) AddResult(creq *page.CurlRequest, result *string) {
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
