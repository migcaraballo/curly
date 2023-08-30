package core

import (
	"fmt"
	"time"
)

type CurlHistoryItem struct {
	Request    *CurlRequest
	CurlResult *string
}

func (c CurlHistoryItem) GetFormattedDate() string {
	return fmt.Sprintf(fmt.Sprintf("%s", c.Request.RequestDate.Format(time.RFC822Z)))
}
