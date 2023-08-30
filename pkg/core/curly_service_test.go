package core

import (
	"fmt"
	"testing"
)

func TestCurlyService_curlIt(t *testing.T) {
	cs := NewCurlyService()
	cs.CheckCurl()

	creq := NewCurlRequest()
	creq.TlsVer = "1.2"
	creq.Method = "GET"
	creq.Url = "https://api.dictionaryapi.dev/api/v2/entries/en/test"

	res, err := cs.ExecuteCurlCall(creq)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}
