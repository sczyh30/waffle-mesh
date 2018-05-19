package route

import (
	"net/http"
	"github.com/sczyh30/waffle-mesh/api/gen"
	"errors"
	"fmt"
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
		r.handleError(writer, err, 404)
		return STOP
	}

	routeAction, err := r.findFirstMatchingRouteAction(routes, request)
	if err != nil {
		r.handleError(writer, err, 404)
		return STOP
	}

	r.executeRouteAction(routeAction, writer, request)

	return STOP
}

func (r *Router) executeRouteAction(action *api.RouteAction, writer http.ResponseWriter, request *http.Request) error {

	return nil
}

func (r *Router) handleError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	fmt.Fprint(w, err.Error())
}

func (r *Router) findFirstMatchingRouteAction(routes []*api.RouteEntry, request *http.Request) (*api.RouteAction, error) {
	path := request.URL.Path

	for _, curRoute := range routes {
		// First match the path pattern.
		if r.matchPathPattern(curRoute, path) {
			headerMatches := curRoute.Match.Headers
			// Then match the header pattern.
			if headerMatches == nil || r.matchHeaderPattern(headerMatches, request.Header) {
				return curRoute.GetRoute(), nil
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

