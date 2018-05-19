package discovery

import (
	"github.com/sczyh30/waffle-mesh/api/gen"
	"github.com/sczyh30/waffle-mesh/brain/k8s/crd"
	"golang.org/x/net/context"
)

type RouteDiscoveryServiceImpl struct {
	Controller *crd.RouteRuleController
}

func (s *RouteDiscoveryServiceImpl) RetrieveRoutes(c context.Context, req *api.DiscoveryRequest) (*api.RouteDiscoveryResponse, error) {
	return &api.RouteDiscoveryResponse{

	}, nil
}
