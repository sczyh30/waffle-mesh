package server

import (
	"fmt"
	"log"

	"github.com/sczyh30/waffle-mesh/proxy/metrics"
	"github.com/sczyh30/waffle-mesh/proxy/discovery"
	"github.com/sczyh30/waffle-mesh/proxy/network"
	"github.com/sczyh30/waffle-mesh/proxy/network/config"
	"github.com/sczyh30/waffle-mesh/proxy/route"
)

const(
	DefaultBrainServerHost = "waffle-brain"
	DefaultGrpcPort = 24242
	DefaultListenerPort = 9080
	DefaultMetricsPort = 19802
)

type ProxyArgs struct {
	BrainServerHost string
	GrpcPort uint32
	ListenerPort uint32
	MetricsPort uint32
}

type startHandler func(chan struct{}) error

type ProxyServer struct {
	metricsServer *metrics.MonitorServer
	xdsWatcher *discovery.ResourceWatcher
	listener *network.Listener

	startHandlerChain []startHandler
}

func NewProxy(args ProxyArgs) (*ProxyServer, error) {
	proxy := &ProxyServer{}

	if err := proxy.initMetricsServer(&args); err != nil {
		return nil, err
	}
	if err := proxy.initDiscoveryWatcher(&args); err != nil {
		return nil, err
	}
	if err := proxy.initListener(&args); err != nil {
		return nil, err
	}

	return proxy, nil
}

func (s *ProxyServer) AddStartHandler(h startHandler) {
	s.startHandlerChain = append(s.startHandlerChain, h)
}

func (s *ProxyServer) StartProxy(stop chan struct{}) error {
	log.Println("Initializing the Waffle Proxy server...")
	for _, fn := range s.startHandlerChain {
		if err := fn(stop); err != nil {
			return err
		}
	}

	return nil
}

func (s *ProxyServer) initMetricsServer(args *ProxyArgs) error {
	// TODO
	return nil
}

func (s *ProxyServer) initDiscoveryWatcher(args *ProxyArgs) error {
	address := args.BrainServerHost + ":" + fmt.Sprint(args.GrpcPort)
	consumer, err := discovery.NewXdsConsumer(address)
	if err != nil {
		return err
	}
	watcher := &discovery.ResourceWatcher{
		XdsConsumer: consumer,
	}

	s.xdsWatcher = watcher
	s.AddStartHandler(func(stop chan struct{}) error {
		go watcher.StartWatching(stop)
		return nil
	})

	return nil
}

func (s *ProxyServer) initListener(args *ProxyArgs) error {
	listener := network.NewListener(network.HTTP2, config.ServerConfig{
		Host: "localhost",
		Port: 8080,
		TlsConfig: config.TlsConfig{
			CertFilePath: "/Users/sczyh30/dev/go-projects/src/github.com/sczyh30/waffle-mesh/cert/cert.pem",
			KeyFilePath:  "/Users/sczyh30/dev/go-projects/src/github.com/sczyh30/waffle-mesh/cert/key.pem",
		},
	})

	s.listener = &listener

	// Configure handlers.
	listener.AddHandler(&route.Router{})

	s.AddStartHandler(func(stop chan struct{}) error {
		go listener.BindAndListen()

		return nil
	})

	return nil
}
