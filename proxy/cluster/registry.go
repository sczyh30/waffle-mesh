package cluster

import "github.com/sczyh30/waffle-mesh/api/gen"

type ClusterRegistry struct {
	registryMap map[string]*ClusterEntry
}

var runtimeClusterRegistry = ClusterRegistry{
	registryMap: make(map[string]*ClusterEntry),
}

func GetCluster(name string) *ClusterEntry {
	return runtimeClusterRegistry.registryMap[name]
}

func RemoveCluster(name string) *ClusterEntry {
	c := runtimeClusterRegistry.registryMap[name]
	delete(runtimeClusterRegistry.registryMap, name)
	return c
}

func AddCluster(name string, cluster *ClusterEntry) {
	runtimeClusterRegistry.registryMap[name] = cluster
}

func DoUpdate(clusters []*api.Cluster, endpoints []*api.ClusterEndpoints) {
	for _, cluster := range clusters {
		if runtimeClusterRegistry.registryMap[cluster.Name] == nil {

		}
	}
}
