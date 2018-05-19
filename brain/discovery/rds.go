package discovery

import (
	"github.com/sczyh30/waffle-mesh/api/gen"
	"golang.org/x/net/context"
	"github.com/sczyh30/waffle-mesh/brain/k8s"
)

type RouteDiscoveryServiceImpl struct {
	controller *k8s.Controller
}

func (s *RouteDiscoveryServiceImpl) RetrieveRoutes(c context.Context, req *api.DiscoveryRequest) (*api.RouteDiscoveryResponse, error) {
	return &api.RouteDiscoveryResponse{

	}, nil
}
