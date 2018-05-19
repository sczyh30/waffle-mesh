package discovery

import (
	"github.com/sczyh30/waffle-mesh/api/gen"
	"golang.org/x/net/context"
	"github.com/sczyh30/waffle-mesh/brain/k8s"
)

type ClusterDiscoveryServiceImpl struct {
	controller *k8s.Controller
}

func (s *ClusterDiscoveryServiceImpl) RetrieveClusters(c context.Context, req *api.DiscoveryRequest) (*api.ClusterDiscoveryResponse, error) {
	return &api.ClusterDiscoveryResponse{

	}, nil
}
