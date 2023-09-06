package core

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGoCurlRequest(t *testing.T) {
	gcc := NewGoCurlClient()

	req := NewCurlRequest()
	req.Url = "http://www.google.com?key=1&key2=lkjafklsjfs"
	req.Method = http.MethodGet
	req.TlsVer = "1.1"

	req.Headers = []string{
		"Content-Type:application/json",
		"Accept:text/html;text/json;application/json",
	}

	req.QsParams = []string{
		"key1=val1",
		"key2=val2",
		"key3=val3",
	}

	res, err := gcc.CurlCall(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
