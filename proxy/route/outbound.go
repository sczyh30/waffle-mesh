package route

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sczyh30/waffle-mesh/api/gen"
	"github.com/sczyh30/waffle-mesh/proxy/cluster"
)

type OutboundRouter struct {
	matcher *Matcher
}

func NewOutboundRouter() *OutboundRouter {
	return &OutboundRouter{
		matcher: &Matcher{},
	}
}

func (r *OutboundRouter) HandleRequest(writer http.ResponseWriter, request *http.Request) bool {
	host := request.Host

	routes, namespace, err := findMatchingRoutes(host)
	if err != nil {
		r.handleError(writer, err, http.StatusNotFound)
		return StopChain
	}

	routeAction, err := r.findFirstMatchingRouteAction(namespace, routes, request)
	if err != nil {
		r.handleError(writer, err, http.StatusNotFound)
		return StopChain
	}

	r.executeRouteAction(routeAction, writer, request)

	return StopChain
}

func (r *OutboundRouter) executeRouteAction(action *RouteActionWrapper, w http.ResponseWriter, request *http.Request) {
	clusterName := action.clusterPicker.NextCluster()
	targetCluster, exists := cluster.GetCluster(clusterName)
	if !exists {
		r.handleError(w, errors.New("no matching cluster"), http.StatusNotFound)
		return
	}

	client, address, err := targetCluster.NextClient(r.lbMetadata(request))
	if err != nil {
		r.handleError(w, err, http.StatusInternalServerError)
		return
	}

	log.Printf("[Outbound] Cluster name: %s, picked endpoint: %s:%d\n", targetCluster.Name(), address.Host, address.Port)

	targetUrl := "http://" + address.Host + ":" + fmt.Sprint(address.Port) + request.RequestURI
	newRequest, err := http.NewRequest(request.Method, targetUrl, request.Body)
	h := request.Header
	h.Add("x-waffle-proxy-from", request.Host)
	newRequest.Header = h
	response, err := client.Do(newRequest)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "service unavailable: " + err.Error())
		return
	}

	// Write headers.
	for k, header := range response.Header {
		for _, v := range header {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(response.StatusCode)
	// Write response body.
	io.Copy(w, response.Body)
}

func (r *OutboundRouter) lbMetadata(request *http.Request) *cluster.LbMetadata {
	return &cluster.LbMetadata{}
}

func (r *OutboundRouter) handleError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	fmt.Fprint(w, err.Error())
}

func (r *OutboundRouter) findFirstMatchingRouteAction(routeNamespace string, routes []*api.RouteEntry, request *http.Request) (*RouteActionWrapper, error) {
	path := request.URL.Path

	for _, curRoute := range routes {
		// Empty match can work.
		if curRoute.Match == nil {
			return getRouteAction(routeNamespace, curRoute.GetRoute()), nil
		}
		// First match the path pattern.
		if curRoute.Match.PathPattern == nil || r.matchPathPattern(curRoute, path) {
			headerMatches := curRoute.Match.Headers
			// Then match the header pattern.
			if headerMatches == nil || r.matchHeaderPattern(headerMatches, request.Header) {
				return getRouteAction(routeNamespace, curRoute.GetRoute()), nil
			}
		}
	}
	return nil, errors.New("no matching routes")
}

func (r *OutboundRouter) matchPathPattern(curRoute *api.RouteEntry, path string) bool {
	return r.matcher.matchExactPath(curRoute, path) ||
		r.matcher.matchPrefixPath(curRoute, path) || r.matcher.matchRegexPath(curRoute, path)
}

func (r *OutboundRouter) matchHeaderPattern(headerMatches []*api.HeaderMatch, header http.Header) bool {
	for _, m := range headerMatches {
		v := header.Get(m.Name)
		if v == "" || !(r.matcher.matchExactHeader(m, v) || r.matcher.matchRegexHeader(m, v)) {
			return false
		}
	}
	return true
}