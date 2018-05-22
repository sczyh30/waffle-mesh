package route

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/sczyh30/waffle-mesh/api/gen"
	"github.com/sczyh30/waffle-mesh/proxy/cluster"
)

const (
	CONTINUE = true
	STOP = false
)

type Router struct {

}

func (r *Router) HandleRequest(writer http.ResponseWriter, request *http.Request) bool {
	host := request.Host

	routes, err := FindMatchingRoutes(host)
	if err != nil {
		r.handleError(writer, err, http.StatusNotFound)
		return STOP
	}

	routeAction, err := r.findFirstMatchingRouteAction(routes, request)
	if err != nil {
		r.handleError(writer, err, http.StatusNotFound)
		return STOP
	}

	r.executeRouteAction(routeAction, writer, request)

	return STOP
}

func (r *Router) executeRouteAction(action *RouteActionWrapper, w http.ResponseWriter, request *http.Request) {
	clusterName := action.clusterPicker.NextCluster()
	targetCluster := cluster.GetCluster(clusterName)
	if targetCluster == nil {
		r.handleError(w, errors.New("no matching cluster"), http.StatusNotFound)
		return
	}

	client, address, err := targetCluster.NextClient(r.lbMetadata(request))
	if err != nil {
		r.handleError(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Cluster name: %s\n", targetCluster.Name())
	fmt.Fprintf(w, "Picked endpoint: %s:%d\n", address.Host, address.Port)
	fmt.Fprintln(w)

	targetUrl := "https://" + address.Host + ":" + fmt.Sprint(address.Port) + request.RequestURI
	newRequest, err := http.NewRequest(request.Method, targetUrl, request.Body)
	newRequest.Header = request.Header
	response, err := client.Do(newRequest)
	if err != nil {
		w.WriteHeader(503)
		fmt.Fprint(w, "service unavailable: " + err.Error())
		return
	}

	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)
}

func (r *Router) lbMetadata(request *http.Request) *cluster.LbMetadata {
	return &cluster.LbMetadata{}
}

func (r *Router) handleError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	fmt.Fprint(w, err.Error())
}

func (r *Router) findFirstMatchingRouteAction(routes []*api.RouteEntry, request *http.Request) (*RouteActionWrapper, error) {
	path := request.URL.Path

	for _, curRoute := range routes {
		// First match the path pattern.
		if r.matchPathPattern(curRoute, path) {
			headerMatches := curRoute.Match.Headers
			// Then match the header pattern.
			if headerMatches == nil || r.matchHeaderPattern(headerMatches, request.Header) {
				if action, ok := curRoute.Action.(*RouteActionWrapper); ok {
					return action, nil
				} else {
					// Directly convert (CHECK HERE!)
					action := FromAction(curRoute.GetRoute())
					curRoute.Action = action
					return action, nil
				}
			}
		}
	}
	return nil, errors.New("no matching routes")
}

func (r *Router) matchPathPattern(curRoute *api.RouteEntry, path string) bool {
	return matchExactPath(curRoute, path) || matchPrefixPath(curRoute, path) || matchRegexPath(curRoute, path)
}

func (r *Router) matchHeaderPattern(headerMatches []*api.HeaderMatch, header http.Header) bool {
	for _, m := range headerMatches {
		v := header.Get(m.Name)
		if v == "" || !(matchExactHeader(m, v) || matchRegexHeader(m, v)) {
			return false
		}
	}
	return true
}

