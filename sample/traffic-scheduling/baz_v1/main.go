package main

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"fmt"
	"log"
	"os"
	"github.com/sczyh30/waffle-mesh/sample/traffic-scheduling/pkg"
)

var bazV1Port = 5763

func main() {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/").To(handleRequest))

	restful.Add(ws)
	err := http.ListenAndServe(fmt.Sprintf(":%d", bazV1Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleRequest(request *restful.Request, response *restful.Response) {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown host (error)"
	}
	response.WriteEntity(pkg.BazResponse{
		Version: "v1",
		Host: host,
	})
}
