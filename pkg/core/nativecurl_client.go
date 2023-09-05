package core

import "os/exec"

type NativeCurlClient struct {
	useableCurlPath string
}

func NewNativeCurlClient(pathToCurl string) *NativeCurlClient {
	return &NativeCurlClient{
		useableCurlPath: pathToCurl,
	}
}

func (nc NativeCurlClient) CurlCall(cr *CurlRequest) (string, error) {
	vArgs := []string{}
	vArgs = append(vArgs, "-v")
	vArgs = append(vArgs, cr.GetTlsArgs()...)
	vArgs = append(vArgs, cr.GetMethodArgs()...)

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

	cmd := exec.Command(nc.useableCurlPath, vArgs...)
	res := cmd.String() + "\n\n"

	stdout, err := cmd.CombinedOutput()
	res += string(stdout)
	return res, err
}

func (nc NativeCurlClient) CurlType() string {
	return "native"
}
