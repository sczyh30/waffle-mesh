package cluster

import "github.com/sczyh30/waffle-mesh/api/gen"

type LoadBalancer interface {
	PickHost(m *LbMetadata) (*api.HttpAddress, error)
}

type LbMetadata struct {
	HashKey string
}

type EndpointWeightPair struct {
	endpoint *api.Endpoint
	effectiveWeight uint32
	currentWeight uint32
}

type RoundRobinLoadBalancer struct {
	endpoints []*EndpointWeightPair
}

func (lb *RoundRobinLoadBalancer) PickHost(m *LbMetadata) (*api.HttpAddress, error) {
	// A smooth load balancing algorithm for weighted round-robin.
	var total uint32 = 0
	for _, pair := range lb.endpoints {
		pair.currentWeight += pair.effectiveWeight
		total += pair.effectiveWeight
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

func NewSmoothWeightedRoundRobinLoadBalancer(e *api.ClusterEndpoints) *RoundRobinLoadBalancer {
	endpoints := make([]*EndpointWeightPair, len(e.Endpoints))
	for i, endpoint := range e.Endpoints {
		endpoints[i] = &EndpointWeightPair{
			endpoint: endpoint,
			effectiveWeight: endpoint.LbWeight,
			currentWeight: 0,
		}
	}
	return &RoundRobinLoadBalancer{
		endpoints: endpoints,
	}
}

type RandomLoadBalancer struct {

}

func (lb *RandomLoadBalancer) PickHost(m *LbMetadata) (*api.HttpAddress, error) {
	return nil, nil
}

type ConsistentHashLoadBalancer struct {

}

func (lb *ConsistentHashLoadBalancer) PickHost(m *LbMetadata) (*api.HttpAddress, error) {
	return nil, nil
}
