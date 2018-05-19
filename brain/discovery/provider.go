package discovery

import "github.com/sczyh30/waffle-mesh/api/gen"

type DiscoveryProvider struct {
	cds *api.ClusterDiscoveryServiceServer
	eds *api.EndpointDiscoveryServiceServer
	rds *api.RouteDiscoveryServiceServer
}