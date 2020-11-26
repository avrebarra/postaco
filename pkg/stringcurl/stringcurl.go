package stringcurl

import (
	"net/http"

	"moul.io/http2curl"
)

type Scurl string

func FromRequest(req *http.Request) (s Scurl, err error) {
	if req == nil {
		return
	}

	curlstr, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return
	}

	s = Scurl(curlstr.String())

	return
}
