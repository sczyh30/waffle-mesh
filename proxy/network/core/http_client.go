package core

import (
	"net/http"
	"golang.org/x/net/http2"
	"crypto/tls"
)

var publicHttp2Tr = &http2.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

func NewHttp2Client() *http.Client {
	tr := publicHttp2Tr

	httpClient := new(http.Client)
	httpClient.Transport = tr
	return httpClient
}