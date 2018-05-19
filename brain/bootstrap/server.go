package bootstrap

import (
	"log"

	"github.com/sczyh30/waffle-mesh/brain/discovery"
	"github.com/sczyh30/waffle-mesh/brain/k8s"
	"k8s.io/client-go/kubernetes"
)

const (
	DefaultXdsProviderPort = 24242
)

type BrainArgs struct {
	XdsProviderPort uint32
}

type startHandler func(chan struct{}) error

type BrainServer struct {
	discoveryProvider *discovery.DiscoveryProvider

	startHandlerChain []startHandler

	k8sClient kubernetes.Interface
	k8sController *k8s.Controller
}

func NewServer(args BrainArgs) (*BrainServer, error) {
	server := &BrainServer{}

	if err := server.initKubernetesClient(&args); err != nil {
		return nil, err
	}
	if err := server.initKubernetesController(&args); err != nil {
		return nil, err
	}
	if err := server.initDiscoveryProvider(&args); err != nil {
		return nil, err
	}
	if err := server.initMetricsServer(&args); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *BrainServer) AddStartHandler(h startHandler) {
	s.startHandlerChain = append(s.startHandlerChain, h)
}

func (s *BrainServer) Start(stop chan struct{}) error {
	log.Println("Initializing the Waffle Brain server...")
	for _, fn := range s.startHandlerChain {
		if err := fn(stop); err != nil {
			return err
		}
	}
	log.Println("Waffle Brain server started")

	return nil
}

func (s *BrainServer) initKubernetesClient(args *BrainArgs) error {
	return nil
}

func (s *BrainServer) initKubernetesController(args *BrainArgs) error {
	return nil
}

func (s *BrainServer) initDiscoveryProvider(args *BrainArgs) error {
	return nil
}

func (s *BrainServer) initMetricsServer(args *BrainArgs) error {
	return nil
}
