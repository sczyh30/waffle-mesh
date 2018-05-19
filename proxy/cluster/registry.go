package cluster

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
