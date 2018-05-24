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

const (
	DefaultBrainServerHost      = "waffle-brain"
	DefaultGrpcPort             = 24242
	DefaultInboundListenerPort  = 9081
	DefaultOutboundListenerPort = 9080
	DefaultMetricsPort          = 19802
)

type ProxyArgs struct {
	BrainServerHost      string
	BrainGrpcPort        uint32
	InboundListenerPort  uint32
	OutboundListenerPort uint32
	MetricsPort          uint32
}

type startHandler func(chan struct{}) error

type ProxyServer struct {
	metricsServer *metrics.MonitorServer
	xdsWatcher    *discovery.ResourceWatcher

	inboundListener  *network.Listener
	outboundListener *network.Listener

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
	log.Println("Waffle Proxy server started")

	return nil
}

func (s *ProxyServer) initMetricsServer(args *ProxyArgs) error {
	ms := metrics.NewMetricsServer(args.MetricsPort)
	s.metricsServer = &ms

	s.AddStartHandler(func(stop chan struct{}) error {
		go ms.Start(stop)
		return nil
	})
	return nil
}

func (s *ProxyServer) initDiscoveryWatcher(args *ProxyArgs) error {
	address := args.BrainServerHost + ":" + fmt.Sprint(args.BrainGrpcPort)
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
	// Init inbound listener (outside services -> Waffle proxy)
	inboundListener := network.NewListener(network.HTTP1_1, config.ServerConfig{
		Host: "localhost",
		Port: int(args.InboundListenerPort),
	})
	inboundListener.AddHandler(route.NewInboundRouter())
	s.inboundListener = &inboundListener

	// Init outbound listener (local service -> Waffle proxy)
	outboundListener := network.NewListener(network.HTTP2, config.ServerConfig{
		Host: "localhost",
		Port: int(args.OutboundListenerPort),
		TlsConfig: config.TlsConfig{
			CertFilePath: "cert/cert.pem",
			KeyFilePath:  "cert/key.pem",
		},
	})
	outboundListener.AddHandler(route.NewOutboundRouter())
	s.outboundListener = &outboundListener

	s.AddStartHandler(func(stop chan struct{}) error {
		go inboundListener.BindAndListen()
		go outboundListener.BindAndListen()

		return nil
	})

	return nil
}
