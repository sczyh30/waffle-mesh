package cluster

import (
	"reflect"
	"testing"

	"github.com/sczyh30/waffle-mesh/api/gen"
)

func TestRoundRobinLoadBalancerPickHost(t *testing.T) {
	addr1 := api.HttpAddress{Host: "10.4.1.101", Port: 8878}
	addr2 := api.HttpAddress{Host: "10.4.1.102", Port: 8878}
	addr3 := api.HttpAddress{Host: "10.4.1.103", Port: 8879}
	endpoints := &api.ClusterEndpoints{
		Endpoints: []*api.Endpoint{
			{Address: &addr1, LbWeight: 2},
			{Address: &addr2, LbWeight: 8},
			{Address: &addr3, LbWeight: 2},
		},
	}
	lb := NewSmoothWeightedRoundRobinLoadBalancer(endpoints)

	assertThisEqual(lb, &addr2, t)
	assertThisEqual(lb, &addr1, t)
	assertThisEqual(lb, &addr2, t)
	assertThisEqual(lb, &addr2, t)
	assertThisEqual(lb, &addr3, t)
	assertThisEqual(lb, &addr2, t)
	assertThisEqual(lb, &addr2, t)
	assertThisEqual(lb, &addr1, t)
	assertThisEqual(lb, &addr2, t)
	assertThisEqual(lb, &addr2, t)
	assertThisEqual(lb, &addr3, t)
	assertThisEqual(lb, &addr2, t)
}

func TestRoundRobinLoadBalancerPickHost2(t *testing.T) {
	addr1 := api.HttpAddress{Host: "10.4.1.101", Port: 8878}
	addr2 := api.HttpAddress{Host: "10.4.1.102", Port: 8878}
	endpoints := &api.ClusterEndpoints{
		Endpoints: []*api.Endpoint{
			{Address: &addr1, LbWeight: 1},
			{Address: &addr2, LbWeight: 1},
		},
	}
	lb := NewSmoothWeightedRoundRobinLoadBalancer(endpoints)

	assertThisEqual(lb, &addr1, t)
	assertThisEqual(lb, &addr2, t)
	assertThisEqual(lb, &addr1, t)
	assertThisEqual(lb, &addr2, t)
	assertThisEqual(lb, &addr1, t)
	assertThisEqual(lb, &addr2, t)
}

func assertThisEqual(lb LoadBalancer, expected *api.HttpAddress, t *testing.T) {
	actual, _ := lb.PickHost(nil)
	// fmt.Printf("Expected: %v, actual: %v\n", expected, actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("the load balancer chose wrong node! Expected: %v, but actual: %v\n", expected, actual)
	}
}