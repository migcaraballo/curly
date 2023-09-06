package core

import (
	"net/http"
	"testing"
)

func TestNativeCurlClient_CurlCall(t *testing.T) {
	cc := NewNativeCurlClient("/usr/bin/curl")
	req := NewCurlRequest()
	req.Url = "https://www.google.com"
	req.TlsVer = "1.1"
	req.Method = http.MethodGet

	res, err := cc.CurlCall(req)
	if err != nil {
		t.Log(err)
	}

	t.Log(res)
}
