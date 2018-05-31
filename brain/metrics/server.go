package metrics

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"

	"github.com/sczyh30/waffle-mesh/api/gen"
)

type BrainMetricsServer struct {
	ws   *restful.WebService
	port uint32
}

func (s *BrainMetricsServer) Start(stop chan struct{}) error {
	restful.Add(s.ws)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func NewBrainMetricsServer(port uint32) *BrainMetricsServer {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/version").To(handleVersion))
	return &BrainMetricsServer{
		ws:   ws,
		port: port,
	}
}

func handleVersion(request *restful.Request, response *restful.Response) {
	response.WriteEntity(api.VersionInfo{ReleaseVersion: "1.0"})
}
