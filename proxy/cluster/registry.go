package cluster

import (
	"github.com/sczyh30/waffle-mesh/api/gen"
	"github.com/sczyh30/waffle-mesh/proxy/network/core"
	"golang.org/x/time/rate"
)

type ClusterRegistry struct {
	registryMap map[string]*ClusterEntry
}

var runtimeClusterRegistry = ClusterRegistry{
	registryMap: make(map[string]*ClusterEntry),
}

func GetCluster(name string) (*ClusterEntry, bool) {
	e, exists := runtimeClusterRegistry.registryMap[name]
	return e, exists
}

func removeCluster(name string) *ClusterEntry {
	c := runtimeClusterRegistry.registryMap[name]
	delete(runtimeClusterRegistry.registryMap, name)
	return c
}

func addClusterEntry(name string, cluster *ClusterEntry) {
	runtimeClusterRegistry.registryMap[name] = cluster
}

func DoUpdate(clusters []*api.Cluster, endpoints []*api.ClusterEndpoints) {
	// Remove unused clusters.
	maps := make(map[string]bool)
	for _, v := range clusters {
		maps[v.Name] = true
	}
	for _, v := range runtimeClusterRegistry.registryMap {
		if _, exists := maps[v.name]; !exists {
			removeCluster(v.name)
		}
	}
	// Build endpoint map to avoid iteration.
	endpointMap := make(map[string]*api.ClusterEndpoints)
	for _, ep := range endpoints {
		endpointMap[ep.ClusterName] = ep
	}
	// Update each cluster.
	for _, cluster := range clusters {
		if oldCluster, exists := runtimeClusterRegistry.registryMap[cluster.Name]; !exists {
			addNewClusterInternal(cluster, endpointMap[cluster.Name])
		} else {
			updateOldClusterInternal(oldCluster, cluster, endpointMap[cluster.Name])
		}
	}
}

func updateOldClusterInternal(oldEntry *ClusterEntry, clusterConfig *api.Cluster, endpoints *api.ClusterEndpoints) {
	oldEntry.doUpdate(clusterConfig, endpoints)
}

func addNewClusterInternal(cluster *api.Cluster, endpoints *api.ClusterEndpoints) {
	// Build client pool.
	pool := make(ClientPool)
	for _, ep := range endpoints.Endpoints {
		pool[toHostAddress(ep.Address)] = core.NewHttp2Client()
	}
	// Build load balancer.
	lb := newLoadBalancerFrom(cluster, endpoints)
	// Build load balancer.
	entry := &ClusterEntry{
		name: cluster.Name,
		endpoints: endpoints,
		config: cluster,
		clientPool: make(ClientPool),
		rateLimiter: &rate.Limiter{},
		lb: lb,
	}
	addClusterEntry(cluster.Name, entry)
}
