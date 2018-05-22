package cluster

import (
	"net/http"
	"sync"

	"github.com/sczyh30/waffle-mesh/api/gen"
	"golang.org/x/time/rate"
)

type ClusterConfig api.Cluster
type ClientPool map[hostAddress]*http.Client

type hostAddress struct {
	host string
	port uint32
}

type ClusterEntry struct {
	name string
	endpoints *api.ClusterEndpoints
	config *ClusterConfig

	lb LoadBalancer
	clientPool ClientPool
	rateLimiter *rate.Limiter

	mutex *sync.Mutex
}

func (c *ClusterEntry) RateLimiter() *rate.Limiter {
	return c.rateLimiter
}

func (c *ClusterEntry) Name() string {
	return c.name
}

func (c *ClusterEntry) UpdateEndpoints(endpoints *api.ClusterEndpoints) error {
	return nil
}

func (c *ClusterEntry) UpdateClusterConfig(config *ClusterConfig) error {
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
