package runtime

type ClusterMetadataRegistry struct {
	registryMap map[string]string
}

var RuntimeClusterRegistry = ClusterMetadataRegistry{}
