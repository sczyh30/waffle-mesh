package model

import "github.com/sczyh30/waffle-mesh/api/gen"

type ConfigConverter struct {

}

func NewConfigConverter() *ConfigConverter {
	return &ConfigConverter{}
}

func (c *ConfigConverter) BuildOutboundClusters() []*api.Cluster {
	return nil
}

func (c *ConfigConverter) BuildInboundClusters() []*api.Cluster {
	return nil
}

func (c *ConfigConverter) AggregateProxyRouteConfigs() []*api.RouteConfig {
	return nil
}