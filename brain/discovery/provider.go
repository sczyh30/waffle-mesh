package discovery

import "github.com/sczyh30/waffle-mesh/api/gen"

type DiscoveryProvider struct {
	cds *api.ClusterDiscoveryServiceServer
	eds *api.EndpointDiscoveryServiceServer
	rds *api.RouteDiscoveryServiceServer
}

func (p *DiscoveryProvider) Start(stop chan struct{}) error {
	return nil
}