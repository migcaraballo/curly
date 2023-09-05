package core

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"strings"
	"time"
)

/*
curl -vvv -L -H 'Content-Type: application/json' -H 'Accept: application/json' 'http://api01.viewpointcloud.com/v2/mentoroh/bluebeam_session/830-353-339/notification?community=mentoroh&userID=auth0%257C601c83c25a90f50071b98909' -d '{}'
--tlsv1.2 --tls-max 1.2
*/

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

func (cr CurlRequest) GetMethodString() string {
	return fmt.Sprintf("-X %s", cr.Method)
}

func (cr CurlRequest) GetMethodArgs() []string {
	return []string{"-X", cr.Method}
}

func (cr CurlRequest) GetTlsString() string {
	return fmt.Sprintf("--tlsv%s --tls-max %s", cr.TlsVer, cr.TlsVer)
}

func (cr CurlRequest) GetTlsArgs() []string {
	return []string{
		fmt.Sprintf("--tlsv%s", cr.TlsVer),
		"--tls-max",
		cr.TlsVer,
	}
}

func (cr CurlRequest) GetTlsUint() uint16 {
	switch cr.TlsVer {
	case "1.0":
		return tls.VersionTLS10
	case "1.1":
		return tls.VersionTLS11
	case "1.2":
		return tls.VersionTLS12
	case "1.3":
		return tls.VersionTLS11
	default:
		return tls.VersionTLS12
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

func (cr CurlRequest) GetHeadersMap() map[string]string {
	rm := map[string]string{}

	for _, v := range cr.Headers {
		ss := strings.Split(v, ":")
		rm[ss[0]] = ss[1]
	}

	return rm
}

func (cr CurlRequest) GetQueryStringMap() map[string]string {
	rm := map[string]string{}

	for _, v := range cr.QsParams {
		ss := strings.Split(v, "=")
		rm[ss[0]] = ss[1]
	}

	return rm
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
