package route

import (
	"errors"

	"github.com/sczyh30/waffle-mesh/api/gen"
	"reflect"
	"encoding/json"
	"log"
)

type RouteTable map[string]*api.RouteConfig

type RouteActionCache map[string][]*RouteActionWrapper

var routeTable = make(RouteTable)
var routeActionCache = make(RouteActionCache)

func getRouteAction(routeName string, action *api.RouteAction) *RouteActionWrapper {
	ac := routeActionCache[routeName]
	for _, v := range ac {
		if reflect.DeepEqual(v.Action, action) {
			return v
		}
	}
	// No match, new create.
	out := fromAction(action)
	ac = append(ac, out)
	return out
}

func addRouteRule(rule *api.RouteConfig) {
	routeTable[rule.Name] = rule
	var ac []*RouteActionWrapper
	for _, entry := range rule.Routes {
		action := fromAction(entry.GetRoute())
		ac = append(ac, action)
	}
	routeActionCache[rule.Name] = ac
}

func updateRouteRule(newConfig *api.RouteConfig) {
	if oldConfig, exists := routeTable[newConfig.Name]; !exists {
		// Add the new route rule.
		addRouteRule(newConfig)
	} else {
		if !reflect.DeepEqual(newConfig, oldConfig) {
			data, _ := json.Marshal(oldConfig)
			log.Printf("Rule to update (old): %s\n", data)
			data, _ = json.Marshal(newConfig)
			log.Printf("Rule to update (new): %s\n", data)

			routeTable[newConfig.Name] = newConfig
			// TODO Check if changes made.
		}
	}
}

func removeRouteRule(name string) {
	delete(routeTable, name)
	delete(routeActionCache, name)
}

func GetRoutes() []*api.RouteConfig {
	routes := make([]*api.RouteConfig, 0)
	for _, v := range routeTable {
		routes = append(routes, v)
	}
	return routes
}

func DoUpdate(routes []*api.RouteConfig) {
	maps := make(map[string]bool)
	for _, v := range routes {
		maps[v.Name] = true
	}
	// GC deprecated rules
	for _, v := range routeTable {
		if _, exists := maps[v.Name]; !exists {
			removeRouteRule(v.Name)
		}
	}
	// Update new rules.
	for _, newConfig := range routes {
		updateRouteRule(newConfig)
	}
}

func findMatchingRoutes(host string) ([]*api.RouteEntry, string, error) {
	for _, routeConfig := range routeTable {
		for _, domain := range routeConfig.Domains {
			// Match any or match exact domain.
			if domain == "*" || domain == host {
				return routeConfig.Routes, routeConfig.Name, nil
			}
		}
	}
	return nil, "", errors.New("cannot find any matching route rules for target host: " + host)
}

type RouteActionWrapper struct {
	Action *api.RouteAction

	clusterPicker ClusterPicker
}

func fromAction(action *api.RouteAction) *RouteActionWrapper {
	if action.GetCluster() != "" {
		return &RouteActionWrapper{
			Action: action,
			clusterPicker: &SingleClusterPicker{Name: action.GetCluster()},
		}
	} else if action.GetWeightedCluster() != nil {
		return &RouteActionWrapper{
			Action:         action,
			clusterPicker: NewSmoothWeightedClusterPicker(action.GetWeightedCluster()),
		}
	}
	return nil
}
