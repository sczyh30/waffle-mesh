package resource

import (
	"sync"

	"github.com/sczyh30/waffle-mesh/api/gen"
)

type ClusterSelectorPair struct {
	clusterName string
	serviceName string
	selectorMap map[string]string
	weight uint32
}

type XdsResourceRegistry struct {
	ClusterConfigs []*api.Cluster
	ClusterEndpoints []*api.ClusterEndpoints
	RouteRuleConfigs []*api.RouteConfig

	mutex sync.Mutex
}

func (registry *XdsResourceRegistry) updateCache(c []*api.Cluster, e []*api.ClusterEndpoints, r []*api.RouteConfig) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	registry.ClusterConfigs = c
	registry.ClusterEndpoints = e
	registry.RouteRuleConfigs = r
}

var xdsCache = XdsResourceRegistry{
	mutex: sync.Mutex{},
}
