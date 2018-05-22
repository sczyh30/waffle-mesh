package route

import (
	"errors"

	"github.com/sczyh30/waffle-mesh/api/gen"
)

type RouteTable map[string]*api.RouteConfig

var routeTable = make(RouteTable)

func AddRouteRule(rule *api.RouteConfig) {
	routeTable[rule.Name] = rule
	for _, entry := range rule.Routes {
		entry.Action = FromAction(entry.GetRoute())
	}
}

func UpdateRouteRule(newConfig *api.RouteConfig) {
	if oldConfig := routeTable[newConfig.Name]; oldConfig == nil {
		// Add the new route rule.
		AddRouteRule(newConfig)
	} else {
		// Check if changes made.
		for _, entry := range oldConfig.Routes {
			//action := entry.GetAction().(*RouteActionWrapper)

			entry.Action = FromAction(entry.GetRoute())
		}
	}

}

func RemoveRouteRule(name string) {
	delete(routeTable, name)
}

func DoUpdate(routes []*api.RouteConfig) {
	for _, newConfig := range routes {
		UpdateRouteRule(newConfig)
	}
}

func FindMatchingRoutes(host string) ([]*api.RouteEntry, error) {
	for _, routeConfig := range routeTable {
		for _, domain := range routeConfig.Domains {
			// Match any or match exact domain.
			if domain == "*" || domain == host {
				return routeConfig.Routes, nil
			}
		}
	}
	return nil, errors.New("cannot find any matching route rules for target host: " + host)
}

type RouteActionWrapper struct {
	Route *api.RouteAction

	clusterPicker ClusterPicker
}

func FromAction(action *api.RouteAction) *RouteActionWrapper {
	if action.GetCluster() != "" {
		return &RouteActionWrapper{
			Route:         action,
			clusterPicker: &SingleClusterPicker{Name: action.GetCluster()},
		}
	} else if action.GetWeightedCluster() != nil {
		return &RouteActionWrapper{
			Route:         action,
			clusterPicker: NewSmoothWeightedClusterPicker(action.GetWeightedCluster()),
		}
	}
	return nil
}

func (*RouteActionWrapper) isRouteEntry_Action() {}
