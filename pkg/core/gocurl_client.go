package core

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type GoCurlClient struct {
}

func NewGoCurlClient() *GoCurlClient {
	return &GoCurlClient{}
}

func (gc GoCurlClient) CurlCall(cr *CurlRequest) (string, error) {
	tlsv := cr.GetTlsUint()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tlsv,
			MaxVersion: tlsv,
		},
	}

	tr.ForceAttemptHTTP2 = true
	clnt := http.Client{
		Transport: tr,
	}

	// init req obj
	req, err := http.NewRequest(cr.Method, cr.Url, nil)
	if err != nil {
		return "", err
	}

	// parse URL
	if u, e := url.Parse(cr.Url); e != nil {
		return "", e
	} else {
		req.URL = u
	}

	// set method
	req.Method = cr.Method

	// add headers
	if len(cr.Headers) > 0 {
		for k, v := range cr.GetHeadersMap() {
			req.Header.Add(k, v)
		}
	}

	// add querystring params if any
	if len(cr.QsParams) > 0 {
		for k, v := range cr.GetQueryStringMap() {
			req.URL.Query().Add(k, v)
		}
	}

	resp, err := clnt.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	rb, err := httputil.DumpResponse(resp, true)

	results := cr.DebugMessage() + string(rb)
	return results, err
}

func (gc GoCurlClient) CurlType() string {
	return "non-native"
}
