package cluster

import (
	"net/http"
	"sync"

	"github.com/sczyh30/waffle-mesh/api/gen"
	"golang.org/x/time/rate"
	"github.com/sczyh30/waffle-mesh/proxy/network/core"
	"encoding/json"
	"log"
	"reflect"
)

type ClientPool map[hostAddress]*http.Client

type hostAddress struct {
	host string
	port uint32
}

type ClusterEntry struct {
	name string
	endpoints *api.ClusterEndpoints
	config *api.Cluster

	lb LoadBalancer
	clientPool ClientPool
	rateLimiter *rate.Limiter

	mutex sync.RWMutex
}

func (c *ClusterEntry) RateLimiter() *rate.Limiter {
	return c.rateLimiter
}

func (c *ClusterEntry) Name() string {
	return c.name
}

func (c *ClusterEntry) doUpdate(newConfig *api.Cluster, newEndpoints *api.ClusterEndpoints) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	log.Println("==================")
	data, _ := json.Marshal(c.config)
	log.Printf("Old config: %v\n", string(data))
	data, _ = json.Marshal(c.endpoints)
	log.Printf("Old endpoints: %v\n", string(data))
	data, _ = json.Marshal(newConfig)
	log.Printf("New config: %v\n", string(data))
	data, _ = json.Marshal(newEndpoints)
	log.Printf("New endpoints: %v\n", string(data))
	log.Println("==================")

	oldConfig := c.config
	// Update LB.
	if newConfig.LbStrategy != oldConfig.LbStrategy {
		c.lb = newLoadBalancerFrom(newConfig, newEndpoints)
	}
	// Update cluster config.
	c.config = newConfig
	// Update endpoints.
	if !reflect.DeepEqual(c.endpoints, newEndpoints) {
		// Build endpoint map to avoid iteration.
		endpointMap := make(map[hostAddress]bool)
		for _, ep := range newEndpoints.Endpoints {
			endpointMap[toHostAddress(ep.Address)] = true
		}
		// Update client pool.
		for k := range c.clientPool {
			if _, exists := endpointMap[k]; !exists {
				delete(c.clientPool, k)
			}
		}
		for k := range endpointMap {
			if _, exists := c.clientPool[k]; !exists {
				c.clientPool[k] = core.NewHttp2Client()
			}
		}
		// Update endpoints
		c.endpoints = newEndpoints
		// Update LB.
		c.lb.DoModify(newEndpoints.Endpoints)
	}


	return nil
}

func (c *ClusterEntry) NextClient(metadata *LbMetadata) (*http.Client, *api.HttpAddress, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	nextAddress, err := c.lb.PickHost(metadata)
	if err != nil {
		return nil, nil, err
	}
	addr := toHostAddress(nextAddress)
	client := c.clientPool[addr]
	if client == nil {
		client = core.NewHttp2Client()
		c.clientPool[addr] = client
	}
	return client, nextAddress, nil
}

func toHostAddress(address *api.HttpAddress) hostAddress {
	return hostAddress{
		host: address.Host,
		port: address.Port,
	}
}

func newLoadBalancerFrom(config *api.Cluster, endpoints *api.ClusterEndpoints) LoadBalancer {
	switch config.LbStrategy {
	case api.Cluster_ROUND_ROBIN:
		return NewSmoothWeightedRoundRobinLoadBalancer(config.Name, endpoints)
	case api.Cluster_RANDOM:
		return NewRandomLoadBalancer(endpoints)
	default:
		return NewSmoothWeightedRoundRobinLoadBalancer(config.Name, endpoints)
	}
}
