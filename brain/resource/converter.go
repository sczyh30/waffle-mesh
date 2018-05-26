package resource

import (
	"fmt"
	"log"
	"strings"

	"github.com/sczyh30/waffle-mesh/api/gen"
	"github.com/sczyh30/waffle-mesh/brain/k8s"
	"github.com/sczyh30/waffle-mesh/brain/k8s/crd"
	crdV1 "github.com/sczyh30/waffle-mesh/brain/k8s/crd/v1"
)

const (
	OutboundPrefix = "outbound"
	InboundPrefix = "inbound"
	HttpTypePrefix = "http"
)

// TODO: temp constants!
const (
	DefaultClientTimeoutInMs = 5000
)

type ResourceConverter struct {
	k8sController *k8s.Controller
	crdController *crd.RouteRuleController
}

func NewResourceConverter(kc *k8s.Controller, cc *crd.RouteRuleController) *ResourceConverter {
	return &ResourceConverter{
		k8sController: kc,
		crdController: cc,
	}
}

func (c *ResourceConverter) BuildOutboundClusters(routes []*api.RouteConfig, selectors map[string]ClusterSelectorPair) []*api.Cluster {
	var clusters []*api.Cluster
	for _, csPair := range selectors {
		cluster := &api.Cluster{
			Name: csPair.clusterName,
			ConnectTimeout: DefaultClientTimeoutInMs,
			LbStrategy: api.Cluster_ROUND_ROBIN,
		}
		clusters = append(clusters, cluster)
	}
	return clusters
}

func (c *ResourceConverter) BuildClusterEndpoints(selectorMap map[string]ClusterSelectorPair) []*api.ClusterEndpoints {
	var clusterEndpoints []*api.ClusterEndpoints
	for _, s := range selectorMap {
		ce, err := c.buildEndpointForCluster(s)
		if err != nil {
			log.Printf("Error when building endpoints for service <%s> (logical cluster: %s)\n", s.serviceName, s.clusterName)
			continue
		}
		if ce != nil {
			clusterEndpoints = append(clusterEndpoints, ce)
		} else {
			log.Printf("Cannot retrieve service <%s> (logical cluster: %s) when building endpoints\n", s.serviceName, s.clusterName)

		}
	}
	return clusterEndpoints
}

func (c *ResourceConverter) matchAllLabels(expected map[string]string, actual map[string]string) bool {
	for k, v := range expected {
		if actual[k] != v {
			return false
		}
	}
	return true
}

func (c *ResourceConverter) buildEndpointForCluster(selectorCluster ClusterSelectorPair) (*api.ClusterEndpoints, error) {
	result := &api.ClusterEndpoints{ClusterName: selectorCluster.clusterName}
	var endpoints []*api.Endpoint

	service, exists := c.k8sController.GetServiceByName(selectorCluster.serviceName)
	if !exists {
		return nil, nil
	}

	podCache := c.k8sController.GetPodCache()

	for _, ep := range c.k8sController.GetEndpoints() {
		if ep.Name == selectorCluster.serviceName { // TODO: namespace
			for _, subset := range ep.Subsets {
				for _, ea := range subset.Addresses {
					labels, exists := podCache.LabelsByIP(ea.IP)
					// Compare the labels.
					if !exists || !c.matchAllLabels(selectorCluster.selectorMap, labels) {
						// log.Printf("Not match. Expected: %v, actual: %v", selectorCluster.selectorMap, labels)
						continue
					}

					//pod, exists := podCache.GetPodByIP(ea.IP)
					// TODO: check here
					if exists {
						for _, port := range service.Spec.Ports {
							endpoints = append(endpoints, &api.Endpoint{
								Address: &api.HttpAddress{Host: ea.IP, Port: uint32(port.Port)},
							})
						}
					} else {

					}
				}

			}
			result.Endpoints = endpoints
			return result, nil
		}
	}
	return nil, nil
}

func (c *ResourceConverter) BuildInboundClusters() []*api.Cluster {
	return nil
}

func (c *ResourceConverter) parseWeightSelectorPair(serviceName string, weightSelectorPair crdV1.RouteSelectorWeight) (string, ClusterSelectorPair) {
	labelSelectorDesc := c.parseLabelSelectors(weightSelectorPair.Labels)
	// Generate final cluster name.
	clusterName := fmt.Sprintf("%s|%s|%s|%s", OutboundPrefix, serviceName, HttpTypePrefix, labelSelectorDesc)
	clusterSelectorPair := ClusterSelectorPair{
		clusterName: clusterName,
		serviceName: serviceName,
		weight: weightSelectorPair.Weight,
		selectorMap: weightSelectorPair.Labels,
	}
	return clusterName, clusterSelectorPair
}

func (c *ResourceConverter) AggregateProxyRouteConfigs(rules []*crdV1.RouteRule) ([]*api.RouteConfig, map[string]ClusterSelectorPair) {
	var configs []*api.RouteConfig
	selectors := make(map[string]ClusterSelectorPair)

	// TODO: validate the data structure!
	for _, ruleCrd := range rules {
		// Parse route name prefix
		routeNamePrefix := ruleCrd.ObjectMeta.Name
		spec := ruleCrd.Spec
		// Parse service name.
		serviceName := spec.Destination.Name
		// Convert the route match condition to xDS RouteMatch.
		routeMatch := c.convertRouteMatchCondition(&spec.Match)
		var routeAction *api.RouteAction
		// TODO: single cluster and weighted cluster can be handled together.
		if len(spec.Route) == 1 {
			clusterName, clusterSelectorPair := c.parseWeightSelectorPair(serviceName, spec.Route[0])
			// Update cluster-selector table for cache.
			selectors[clusterName] = clusterSelectorPair
			// Generate route action. Cluster pattern is `single cluster`.
			routeAction = &api.RouteAction{
				ClusterPattern: &api.RouteAction_Cluster{Cluster: clusterName},
				TimeoutMs: DefaultClientTimeoutInMs,
				// TODO: Add retry strategy.
			}
		} else {
			var clusterWeightPairs []*api.WeightedCluster_ClusterWeightPair
			for _, weightSelectorPair := range spec.Route {
				clusterName, clusterSelectorPair := c.parseWeightSelectorPair(serviceName, weightSelectorPair)
				// Update cluster-selector table for cache.
				selectors[clusterName] = clusterSelectorPair
				// Update cluster weight pairs.
				clusterWeightPairs = append(clusterWeightPairs, &api.WeightedCluster_ClusterWeightPair{
					Name: clusterName,
					Weight: weightSelectorPair.Weight,
				})
			}
			// Generate route action. Cluster pattern is `weighted cluster`.
			routeAction = &api.RouteAction{
				ClusterPattern: &api.RouteAction_WeightedCluster{
					WeightedCluster: &api.WeightedCluster{
						Clusters: clusterWeightPairs,
					},
				},
				TimeoutMs: DefaultClientTimeoutInMs,
			}
		}

		// Generate name for RouteConfig.
		routeConfigName := fmt.Sprintf("%s|%s", routeNamePrefix, HttpTypePrefix)
		// Build the single route entry.
		routeEntry := &api.RouteEntry{
			Match: routeMatch,
			Action: &api.RouteEntry_Route{
				Route: routeAction,
			},
		}

		// Generate the final config and add to list.
		routeConfig := &api.RouteConfig{
			Name: routeConfigName,
			Domains: c.parseRouteDomains(serviceName),
			Routes: []*api.RouteEntry{routeEntry},
		}
		configs = append(configs, routeConfig)
	}

	return configs, selectors
}

func(c *ResourceConverter) parseLabelSelectors(labelMap map[string]string) string {
	var labels []string
	for l, v := range labelMap {
		labels = append(labels, fmt.Sprintf("%s=%s", l, v))
	}
	return strings.Join(labels, ",")
}

func(c *ResourceConverter) parseRouteDomains(serviceName string) []string {
	var domains []string
	namespace := "default"
	// Add serviceName as a domain.
	domains = append(domains, serviceName)
	// Add serviceName.namespace as a domain.
	domains = append(domains, fmt.Sprintf("%s.%s", serviceName, namespace))
	service, exists := c.k8sController.GetServiceByName(serviceName)
	if !exists {
		log.Printf("Warning: service <%s> does not exist!\n", serviceName)
		return domains
	}
	// Add serviceName:port as a domain.
	for _, port := range service.Spec.Ports {
		domains = append(domains, fmt.Sprintf("%s:%d", serviceName, port.Port))
		domains = append(domains, fmt.Sprintf("%s.%s:%d", serviceName, namespace, port.Port))
	}

	return domains
}

func (c *ResourceConverter) convertRouteMatchCondition(match *crdV1.RouteMatchCondition) *api.RouteMatch {
	convertedMatch := &api.RouteMatch{}
	uriHeaderName := "uri"
	uriMatch := match.Request.Headers[uriHeaderName]
	// Parse path pattern.
	if isConditionValid(uriMatch) {
		if uriMatch.Exact != "" {
			convertedMatch.PathPattern = &api.RouteMatch_ExactPath{ExactPath: uriMatch.Exact}
		} else if uriMatch.Prefix != "" {
			convertedMatch.PathPattern = &api.RouteMatch_Prefix{Prefix: uriMatch.Prefix}
		} else if uriMatch.Regex != "" {
			convertedMatch.PathPattern = &api.RouteMatch_Regex{Regex: uriMatch.Regex}
		}
	}
	// Parse header conditions.
	var headerMatchList []*api.HeaderMatch
	for key, cond := range match.Request.Headers {
		if key == uriHeaderName || !isConditionValid(cond) {
			continue
		}
		headerMatch := &api.HeaderMatch{
			Name: key,
		}
		if cond.Exact != "" {
			headerMatch.HeaderMatchPattern = &api.HeaderMatch_ExactMatch{ExactMatch: cond.Exact}
		} else if cond.Regex != "" {
			headerMatch.HeaderMatchPattern = &api.HeaderMatch_RegexMatch{RegexMatch: cond.Regex}
		}
		headerMatchList = append(headerMatchList, headerMatch)
	}
	convertedMatch.Headers = headerMatchList
	return convertedMatch
}

func isConditionValid(cond crdV1.StringMatchCondition) bool {
	return cond.Exact != "" || cond.Prefix != "" || cond.Regex != ""
}