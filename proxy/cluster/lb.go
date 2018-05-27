package cluster

import (
	"github.com/sczyh30/waffle-mesh/api/gen"
	"reflect"
	"log"
)

// Not thread-safe.
type LoadBalancer interface {
	PickHost(m *LbMetadata) (*api.HttpAddress, error)

	DoModify(endpoints []*api.Endpoint) bool
}

type LbMetadata struct {
	HashKey string
}

type EndpointWeightPair struct {
	endpoint *api.Endpoint
	effectiveWeight uint32
	currentWeight int32
}

type RoundRobinLoadBalancer struct {
	endpoints []*EndpointWeightPair
	existMap map[hostAddress]*api.Endpoint

	clusterName string
}

func (lb *RoundRobinLoadBalancer) PickHost(m *LbMetadata) (*api.HttpAddress, error) {
	// A smooth load balancing algorithm for weighted round-robin.
	var total int32 = 0
	for _, pair := range lb.endpoints {
		// fmt.Printf("%v:%d\n", pair.endpoint, pair.currentWeight)
		pair.currentWeight += int32(pair.effectiveWeight)
		total += int32(pair.effectiveWeight)
	}
	max := lb.endpoints[0]
	for _, pair := range lb.endpoints {
		if pair.currentWeight > max.currentWeight {
			max = pair
		}
	}

	max.currentWeight -= total

	return max.endpoint.Address, nil
}

func (lb *RoundRobinLoadBalancer) DoModify(endpoints []*api.Endpoint) bool {
	canModify := func() bool {
		for _, newEndpoint := range endpoints {
			if old, exists := lb.existMap[toHostAddress(newEndpoint.Address)]; !exists {
				return true
			} else {
				if !reflect.DeepEqual(old, newEndpoint) {
					return true
				}
			}
		}
		return false
	} ()
	if canModify {
		lb.reset(endpoints)
		log.Printf("Load balancer for cluster <%s> modified\n", lb.clusterName)
	}
	return canModify
}

func (lb *RoundRobinLoadBalancer) reset(e []*api.Endpoint) {
	endpoints, existsMap := fromEndpoints(e)
	lb.endpoints = endpoints
	lb.existMap = existsMap
}

func fromEndpoints(endpoints []*api.Endpoint) ([]*EndpointWeightPair, map[hostAddress]*api.Endpoint) {
	pair := make([]*EndpointWeightPair, len(endpoints))
	existsMap := make(map[hostAddress]*api.Endpoint)
	for i, endpoint := range endpoints {
		existsMap[toHostAddress(endpoint.Address)] = endpoint
		pair[i] = &EndpointWeightPair{
			endpoint: endpoint,
			effectiveWeight: endpoint.LbWeight,
			currentWeight: 0,
		}
	}
	return pair, existsMap
}

func NewSmoothWeightedRoundRobinLoadBalancer(name string, e *api.ClusterEndpoints) *RoundRobinLoadBalancer {
	endpoints, existsMap := fromEndpoints(e.Endpoints)
	return &RoundRobinLoadBalancer{
		clusterName: name,
		endpoints: endpoints,
		existMap: existsMap,
	}
}

type RandomLoadBalancer struct {
	endpoints []*api.Endpoint
}

func (lb *RandomLoadBalancer) PickHost(m *LbMetadata) (*api.HttpAddress, error) {
	return nil, nil
}

func (lb *RandomLoadBalancer) DoModify(endpoints []*api.Endpoint) bool {
	return false
}

func NewRandomLoadBalancer(e *api.ClusterEndpoints) *RandomLoadBalancer {
	return &RandomLoadBalancer{
		endpoints: e.Endpoints,
	}
}

type ConsistentHashLoadBalancer struct {

}

func (lb *ConsistentHashLoadBalancer) PickHost(m *LbMetadata) (*api.HttpAddress, error) {
	return nil, nil
}

func (lb *ConsistentHashLoadBalancer) DoModify(endpoints []*api.Endpoint) bool {
	return false
}
