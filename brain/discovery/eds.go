package discovery

import (
	"github.com/sczyh30/waffle-mesh/api/gen"
	"github.com/sczyh30/waffle-mesh/brain/resource"
	"golang.org/x/net/context"
)

type EndpointDiscoveryServiceImpl struct {
}

func (s *EndpointDiscoveryServiceImpl) RetrieveEndpoints(c context.Context, req *api.DiscoveryRequest) (*api.EndpointDiscoveryResponse, error) {
	return &api.EndpointDiscoveryResponse{
		Result: resource.XdsCache.ClusterEndpoints,
	}, nil
}