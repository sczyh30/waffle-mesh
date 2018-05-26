package network

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/sczyh30/waffle-mesh/proxy/network/config"
	"github.com/sczyh30/waffle-mesh/proxy/network/core"
)

// Proxy listener observes the port and process the requests.
type Listener interface {
	AddHandler(handler HttpHandler)

	BindAndListen() error
}

type ServerType int

const (
	HTTP1_1 ServerType = iota
	HTTP2
)

type listenerImpl struct {
	serverType ServerType
	server http.Server

	handlerChain []HttpHandler
	config config.ServerConfig

	mutex sync.RWMutex
}

func (l *listenerImpl) AddHandler(handler HttpHandler) {
	l.handlerChain = append(l.handlerChain, handler)
}

func (l *listenerImpl) BindAndListen() error {
	if len(l.handlerChain) == 0 {
		return errors.New("empty handler chain")
	}

	var err error

	mux := http.NewServeMux()
	mux.HandleFunc("/", l.handleRequest)
	l.server.Handler = mux

	switch l.serverType {
	case HTTP1_1:
		err = l.server.ListenAndServe()
	case HTTP2:
		err = l.server.ListenAndServe()
		//err = l.server.ListenAndServeTLS(l.config.TlsConfig.CertFilePath, l.config.TlsConfig.KeyFilePath)
	}
	if err != nil {
		log.Fatal("error when listening to port " + strconv.Itoa(l.config.Port), err)
	}
	return err
}

func (l *listenerImpl) handleRequest(w http.ResponseWriter, r *http.Request)  {
	for _, handler := range l.handlerChain {
		handler.HandleRequest(w, r)
	}
}

func NewListener(serverType ServerType, config config.ServerConfig) Listener {
	l := &listenerImpl{
		serverType: serverType,
		config: config,
		mutex: sync.RWMutex{},
	}
	if serverType == HTTP2 {
		l.server = core.NewHttp2Server()
	} else {
		l.server = core.NewHttpServer()
	}
	// Resolve host and port.
	addr := l.config.Host + ":" + strconv.Itoa(l.config.Port)
	l.server.Addr = addr

	return l
}