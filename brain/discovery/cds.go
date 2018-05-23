package discovery

import (
	"github.com/sczyh30/waffle-mesh/api/gen"
	"golang.org/x/net/context"
	"github.com/sczyh30/waffle-mesh/brain/resource"
)

type ClusterDiscoveryServiceImpl struct {
}

func (s *ClusterDiscoveryServiceImpl) RetrieveClusters(c context.Context, req *api.DiscoveryRequest) (*api.ClusterDiscoveryResponse, error) {
	return &api.ClusterDiscoveryResponse{
		Clusters: resource.XdsCache.ClusterConfigs,
	}, nil
}
