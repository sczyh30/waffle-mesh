package cluster

import (
	"net/http"

	"github.com/sczyh30/waffle-mesh/api/gen"
)

type ClusterConfig api.Cluster
type ClientPool map[hostAddress]*http.Client

type hostAddress struct {
	host string
	port uint32
}

type ClusterEntry struct {
	Name string
	endpoints *api.ClusterEndpoints
	config *ClusterConfig

	lb LoadBalancer
	clientPool ClientPool
}

func (c *ClusterEntry) updateEndpoints(endpoints *api.ClusterEndpoints) error {
	return nil
}

func (c *ClusterEntry) updateClusterConfig(config *ClusterConfig) error {
	return nil
}

func (c *ClusterEntry) NextClient(metadata *LbMetadata) (*http.Client, *api.HttpAddress, error) {
	nextAddress, err := c.lb.PickHost(metadata)
	if err != nil {
		return nil, nil, err
	}
	addr := hostAddress{
		host: nextAddress.Host,
		port: nextAddress.Port,
	}
	client := c.clientPool[addr]
	if client == nil {
		// TODO: Init the client pool
	}
	return client, nextAddress, nil
}
