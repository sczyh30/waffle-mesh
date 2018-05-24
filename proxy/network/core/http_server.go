package core

import (
	"net/http"
)

func NewHttp2Server() http.Server {
	var server http.Server
	//http2.ConfigureServer(&server, &http2.Server{})
	return server
}

func NewHttpServer() http.Server {
	var server http.Server
	return server
}