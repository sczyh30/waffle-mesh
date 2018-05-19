package discovery

import (
	"github.com/sczyh30/waffle-mesh/brain/k8s"
	"github.com/sczyh30/waffle-mesh/api/gen"
	"golang.org/x/net/context"
)

type EndpointDiscoveryServiceImpl struct {
	Controller *k8s.Controller
}

func (s *EndpointDiscoveryServiceImpl) RetrieveEndpoints(c context.Context, req *api.DiscoveryRequest) (*api.EndpointDiscoveryResponse, error) {
	return &api.EndpointDiscoveryResponse{

	}, nil
}