package network

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"net/http"
	"strconv"
	"sync"

	"github.com/sczyh30/waffle-mesh/proxy/network/config"
	"github.com/sczyh30/waffle-mesh/proxy/network/core"
	"github.com/sczyh30/waffle-mesh/proxy/route"
	"github.com/sczyh30/waffle-mesh/proxy/runtime"

)

// Proxy listener observes the port and process the requests.
type Listener interface {
	AddHandler(handler *HttpHandler)

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

	handlerChain []*HttpHandler
	config config.ServerConfig

	mutex *sync.RWMutex
}

func (l *listener) AddHandler(handler *HttpHandler) {
	l.handlerChain = append(l.handlerChain, handler)
}

func (l *listener) BindAndListen() error {
	var err error
	// Resolve host and port.
	addr := l.config.Host + ":" + strconv.Itoa(l.config.Port)
	l.server.Addr = addr

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRequest)
	l.server.Handler = mux

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

func handleRequest(w http.ResponseWriter, r *http.Request)  {
	client := core.NewHttp2Client()

	path := r.URL.Path
	host := r.Host
	//clientAddr := r.RemoteAddr
	method := r.Method

	routes, err := route.FindMatchingRoutes(host)
	if err != nil {
		//log.Fatal(err)
		w.WriteHeader(404)
		fmt.Fprint(w, "no matching routes")
		return
	}
	var targetClusterName = ""
	for _, curRoute := range routes {
		if curRoute.Match.GetExactPath() != "" {
			// Match exact path.
			if route.MatchExact(curRoute.Match.GetExactPath(), path) {
				// TODO!
				targetClusterName = curRoute.GetRoute().GetCluster()
			}
		} else if curRoute.Match.GetPrefix() != "" {
			// Match path prefix.
			if route.MatchPrefix(curRoute.Match.GetPrefix(), path) {
				// TODO!
				targetClusterName = curRoute.GetRoute().GetCluster()
			}
		} else if curRoute.Match.GetRegex() != "" {
			if route.MatchRegex(curRoute.Match.GetRegex(), path) {
				// TODO!
				targetClusterName = curRoute.GetRoute().GetCluster()
			}
			// Match regex pattern.
		}
	}
	// No matching
	if targetClusterName == "" {
		w.WriteHeader(404)
		fmt.Fprint(w, "no matching route")
		return
	}
	targetCluster := runtime.GetCluster(targetClusterName)
	if targetCluster == nil {
		w.WriteHeader(404)
		fmt.Fprint(w, "no matching cluster")
		return
	}

	fmt.Fprintf(w, "Cluster name: %s\n", targetCluster.Name)
	fmt.Fprintln(w, "Cluster registered endpoints:")
	for _, address := range targetCluster.Hosts {
		fmt.Fprintf(w, "host: %s, port: %d\n", address.Host, address.Port)
	}
	fmt.Fprintln(w)

	targetUrl, _ := url.Parse("https://" + targetCluster.Hosts[0].Host + ":" + fmt.Sprint(targetCluster.Hosts[0].Port) + r.RequestURI)
	response, err := client.Do(&http.Request{
		Method: method,
		URL: targetUrl,
		Body: r.Body,
		Header: r.Header,
	})
	if err != nil {
		w.WriteHeader(503)
		fmt.Fprint(w, "service unavailable: " + err.Error())
		return
	}
	io.Copy(w, response.Body)
}

func NewListener(serverType ServerType, config config.ServerConfig) Listener {
	l := &listener{
		serverType: serverType,
		config: config,
	}
	if serverType == HTTP2 {
		l.server = core.NewHttp2Server()
	} else {
		l.server = core.NewHttpServer()
	}
	return l
}