package network

import (
	"github.com/sczyh30/waffle-mesh/proxy/network/config"
	"net/http"
	"container/list"
	"strconv"
	"log"
	"github.com/sczyh30/waffle-mesh/proxy/network/core"
)

// Proxy listener observes the port and process the requests.
type Listener interface {
	AddInboundProcessor(processor InboundProcessor) error

	AddOutboundProcessor(processor OutboundProcessor) error

	BindAndListen() error
}

type ServerType int

const (
	HTTP1_1 ServerType = iota
	HTTP2
)

type listener struct {
	serverType ServerType
	server http.Server

	inboundProcessorChain *list.List
	outboundProcessorChain *list.List
	config config.ServerConfig
}

func (l *listener) AddInboundProcessor(processor InboundProcessor) error {
	return nil
}

func (l *listener) AddOutboundProcessor(processor OutboundProcessor) error {
	return nil
}

func (l *listener) BindAndListen() error {
	var err error
	// Resolve host and port.
	addr := l.config.Host + ":" + strconv.Itoa(l.config.Port)
	l.server.Addr = addr

	//mux := http.NewServeMux()

	switch l.serverType {
	case HTTP1_1:
		err = l.server.ListenAndServe()
	case HTTP2:
		err = l.server.ListenAndServeTLS(l.config.TlsConfig.CertFilePath, l.config.TlsConfig.KeyFilePath)
	}
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	return err
}

func NewListener(serverType ServerType, config config.ServerConfig) Listener {
	l := &listener{
		serverType: serverType,
		inboundProcessorChain: list.New(),
		outboundProcessorChain: list.New(),
		config: config,
	}
	if serverType == HTTP2 {
		l.server = core.NewHttp2Server()
	} else {
		l.server = core.NewHttpServer()
	}
	return l
}