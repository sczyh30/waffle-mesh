package discovery

import (
	"context"

	"github.com/sczyh30/waffle-mesh/api/gen"
	"github.com/sczyh30/waffle-mesh/proxy/cluster"
	"github.com/sczyh30/waffle-mesh/proxy/route"
	"google.golang.org/grpc"
)

type XdsConsumer struct {
	cdsClient api.ClusterDiscoveryServiceClient
	edsClient api.EndpointDiscoveryServiceClient
	rdsClient api.RouteDiscoveryServiceClient

	conn *grpc.ClientConn
}

func NewXdsConsumer(address string) (*XdsConsumer, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &XdsConsumer{
		conn: conn,
		cdsClient: api.NewClusterDiscoveryServiceClient(conn),
		edsClient: api.NewEndpointDiscoveryServiceClient(conn),
		rdsClient: api.NewRouteDiscoveryServiceClient(conn),
	}, nil
}

func (c *XdsConsumer) RetrieveAndUpdate(ctx context.Context) error {
	cdsResponse, err := c.cdsClient.RetrieveClusters(ctx, &api.DiscoveryRequest{})
	if err != nil {
		return err
	}
	edsResponse, err := c.edsClient.RetrieveEndpoints(ctx, &api.DiscoveryRequest{})
	if err != nil {
		return err
	}
	rdsResponse, err := c.rdsClient.RetrieveRoutes(ctx, &api.DiscoveryRequest{})
	if err != nil {
		return err
	}

	// Apply changes to cluster registry.
	cluster.DoUpdate(cdsResponse.Clusters, edsResponse.Result)
	// Apply changes to route table.
	route.DoUpdate(rdsResponse.Result)

	return nil
}
