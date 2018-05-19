package cluster

import "github.com/sczyh30/waffle-mesh/api/gen"

type LoadBalancer interface {
	PickHost(m *LbMetadata) (*api.HttpAddress, error)
}

type LbMetadata struct {
	HashKey string
}

type RoundRobinLoadBalancer struct {

}

func (lb *RoundRobinLoadBalancer) PickHost(m *LbMetadata) (*api.HttpAddress, error) {
	return nil, nil
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
