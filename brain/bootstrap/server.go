package bootstrap

import (
	"log"

	"github.com/sczyh30/waffle-mesh/brain/discovery"
	"github.com/sczyh30/waffle-mesh/brain/k8s"
	"k8s.io/client-go/kubernetes"
	"github.com/sczyh30/waffle-mesh/brain/k8s/crd"
	"github.com/sczyh30/waffle-mesh/brain/resource"
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
	k8sRouteRuleController *crd.RouteRuleController

	xdsUpdater *resource.XdsResourceUpdater
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
	_, k8sClient, err := k8s.CreateInClusterInterface()
	s.k8sClient = k8sClient
	return err
}

func (s *BrainServer) initKubernetesController(args *BrainArgs) error {
	options := k8s.ControllerOptions{
		WatchedNamespace: "",
	}
	k8sController := k8s.NewController(s.k8sClient, options)
	s.k8sController = k8sController

	routeRuleController, err := crd.NewRouteRuleController(options)
	if err != nil {
		return err
	}
	s.k8sRouteRuleController = routeRuleController

	updater := resource.NewXdsResourceUpdater(k8sController, routeRuleController)
	s.xdsUpdater = updater

	s.AddStartHandler(func(stop chan struct{}) error {
		go s.k8sController.Run(stop)
		go s.k8sRouteRuleController.Run(stop)

		return nil
	})

	return nil
}

func (s *BrainServer) initDiscoveryProvider(args *BrainArgs) error {
	provider := &discovery.DiscoveryProvider{
		Port: args.XdsProviderPort,
		Cds: &discovery.ClusterDiscoveryServiceImpl{
			Controller: s.k8sController,
		},
		Eds: &discovery.EndpointDiscoveryServiceImpl{
			Controller: s.k8sController,
		},
		Rds: &discovery.RouteDiscoveryServiceImpl{
			Controller: s.k8sRouteRuleController,
		},
	}
	s.discoveryProvider = provider

	s.AddStartHandler(func(stop chan struct{}) error {
		go s.xdsUpdater.Start(stop)

		go s.discoveryProvider.Start(stop)

		return nil
	})

	return nil
}

func (s *BrainServer) initMetricsServer(args *BrainArgs) error {
	return nil
}
