package core

import (
	"bytes"
	"fmt"
	"time"
)

type CurlRequest struct {
	RequestDate time.Time
	Url         string
	TlsVer      string
	Method      string
	Headers     []string
	QsParams    []string
	Body        string
}

func NewCurlRequest() *CurlRequest {
	return &CurlRequest{
		RequestDate: time.Time{},
		Url:         "",
		TlsVer:      "",
		Method:      "",
		Headers:     []string{},
		QsParams:    []string{},
		Body:        "",
	}
}

func (cr CurlRequest) GetBodyString() string {
	if len(cr.Url) > 0 {
		return fmt.Sprintf("-d %s", cr.Body)
	}

	return ""
}

func (cr CurlRequest) GetHeadersString() string {
	b := bytes.Buffer{}

	if len(cr.Headers) > 0 {
		for _, v := range cr.Headers {
			b.WriteString(fmt.Sprintf("-H \"%s\" ", v))
		}
	}

	return b.String()
}

func (cr CurlRequest) GetQsParamString() string {
	b := bytes.Buffer{}

	if len(cr.QsParams) > 0 {
		for _, v := range cr.QsParams {
			b.WriteString(fmt.Sprintf("--data-urlencode \"%s\" ", v))
		}
	}

	return b.String()
}

func (cr CurlRequest) DebugMessage() string {
	d := bytes.Buffer{}

	d.WriteString("Debug Start: ------------------------------\n")
	d.WriteString(fmt.Sprintf("- Url: %s\n", cr.Url))
	d.WriteString(fmt.Sprintf("- Method: %s\n", cr.Method))
	d.WriteString(fmt.Sprintf("- TLS: %s\n", cr.TlsVer))
	d.WriteString(fmt.Sprintf("- Headers %s\n", cr.GetHeadersString()))
	d.WriteString(fmt.Sprintf("- QsParams: %s\n", cr.GetQsParamString()))
	d.WriteString(fmt.Sprintf("- Body: %s\n", cr.GetBodyString()))
	d.WriteString("Debug End: ------------------------------\n")

	return d.String()
}
