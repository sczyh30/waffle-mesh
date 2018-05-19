package cluster

import (
	"net/http"

	"github.com/sczyh30/waffle-mesh/api/gen"
)

type ClusterConfig api.Cluster
type ClientPool map[api.HttpAddress]*http.Client

type ClusterEntry struct {
	name string
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

func (c *ClusterEntry) NextClient(metadata *LbMetadata) (*http.Client, error) {
	nextAddress, err := c.lb.PickHost(metadata)
	if err != nil {
		return nil, err
	}
	client := c.clientPool[*nextAddress]
	if client == nil {
		// TODO: Init the client pool
	}
	return client, nil
}
