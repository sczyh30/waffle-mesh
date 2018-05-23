package discovery

import (
	"github.com/sczyh30/waffle-mesh/api/gen"
	"github.com/sczyh30/waffle-mesh/brain/resource"
	"golang.org/x/net/context"
)

type RouteDiscoveryServiceImpl struct {

}

func (s *RouteDiscoveryServiceImpl) RetrieveRoutes(c context.Context, req *api.DiscoveryRequest) (*api.RouteDiscoveryResponse, error) {
	return &api.RouteDiscoveryResponse{
		Result: resource.XdsCache.RouteRuleConfigs,
	}, nil
}
