package metrics

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/sczyh30/waffle-mesh/proxy/route"
	"github.com/sczyh30/waffle-mesh/proxy/cluster"
	"github.com/sczyh30/waffle-mesh/api/gen"
)

type MonitorServer interface {
	Start(stop chan struct{}) error
}

type simpleMetricsServer struct {
	ws   *restful.WebService
	port uint32
}

func (s *simpleMetricsServer) Start(stop chan struct{}) error {
	restful.Add(s.ws)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func NewMetricsServer(port uint32) MonitorServer {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/version").To(handleVersion))
	ws.Route(ws.GET("/route_rules").To(handleGetRouteRules))
	ws.Route(ws.GET("/clusters").To(handleGetClusters))
	ws.Route(ws.GET("/endpoints").To(handleGetClusterEndpoints))
	return &simpleMetricsServer{
		ws:   ws,
		port: port,
	}
}

func handleVersion(request *restful.Request, response *restful.Response) {
	response.WriteEntity(api.VersionInfo{ReleaseVersion: "1.0"})
}

func handleGetRouteRules(request *restful.Request, response *restful.Response) {
	response.WriteEntity(route.GetRoutes())
}

func handleGetClusters(request *restful.Request, response *restful.Response) {
	response.WriteEntity(cluster.GetClusters())
}

func handleGetClusterEndpoints(request *restful.Request, response *restful.Response) {
	response.WriteEntity(cluster.GetClusterEndpoints())
}
