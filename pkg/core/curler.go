package core

type Curler interface {
	CurlCall(cr *CurlRequest) (string, error)
	CurlType() string
}
