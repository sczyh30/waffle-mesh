package discovery

import (
	"github.com/sczyh30/waffle-mesh/api/gen"
	"net"
	"fmt"
	"google.golang.org/grpc"
)

type DiscoveryProvider struct {
	Port uint32

	GrpcServer *grpc.Server
	Cds api.ClusterDiscoveryServiceServer
	Eds api.EndpointDiscoveryServiceServer
	Rds api.RouteDiscoveryServiceServer
}

func (p *DiscoveryProvider) Start(stop chan struct{}) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", p.Port))
	if err != nil {
		return err
	}
	s := grpc.NewServer()

	api.RegisterClusterDiscoveryServiceServer(s, p.Cds)
	api.RegisterEndpointDiscoveryServiceServer(s, p.Eds)
	api.RegisterRouteDiscoveryServiceServer(s, p.Rds)

	p.GrpcServer = s

	if err := s.Serve(listener); err != nil {
		return err
	}

	return nil
}