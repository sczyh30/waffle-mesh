package route

import (
	"github.com/sczyh30/waffle-mesh/api/gen"
	"errors"
)

type RouteActionWrapper struct {
	Route *api.RouteAction

	clusterPicker ClusterPicker
}

func FromAction(action *api.RouteAction) *RouteActionWrapper {
	if action.GetCluster() != "" {
		return &RouteActionWrapper{
			Route: action,
			clusterPicker: &SingleClusterPicker{name: action.GetCluster()},
		}
	} else if action.GetWeightedCluster() != nil {
		return &RouteActionWrapper{
			Route: action,
			clusterPicker: NewSmoothWeightedClusterPicker(action.GetWeightedCluster()),
		}
	}
	return nil
}

func (*RouteActionWrapper) isRouteEntry_Action() {}

var routeRuleRegistry = make(map[string]*api.RouteConfig)

func AddRouteRule(name string, rule *api.RouteConfig) {
	routeRuleRegistry[name] = rule
	for _, details := range rule.RouteDetails {
		for _, entry := range details.Routes {
			entry.Action = FromAction(entry.GetRoute())
		}
	}
}

func RemoveRouteRule(name string) {
	delete(routeRuleRegistry, name)
}

func FindMatchingRoutes(host string) ([]*api.RouteEntry, error) {
	for _, v := range routeRuleRegistry {
		for _, rd := range v.RouteDetails {
			for _, domain := range rd.Domains {
				// Match any or match exact domain.
				if domain == "*" || domain == host {
					return rd.Routes, nil
				}
			}
		}
	}
	return nil, errors.New("cannot find any matching route rules for target host: " + host)
}
