package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"

	"github.com/emicklei/go-restful"
	"github.com/sczyh30/waffle-mesh/sample/traffic-scheduling/pkg"
)

var fooPort = 5758

func main() {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/baz").To(handleBazInvocation))

	restful.Add(ws)
	err := http.ListenAndServe(fmt.Sprintf(":%d", fooPort), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleBazInvocation(request *restful.Request, response *restful.Response) {
	client := http.Client{}
	resp, err := client.Get("http://baz")
	if err != nil {
		response.WriteError(http.StatusServiceUnavailable, err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		response.WriteError(resp.StatusCode, errors.New(string(body)))
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	bazResp := pkg.BazResponse{}
	err = json.Unmarshal(body, &bazResp)

	data, _ := json.Marshal(pkg.FooResponse{Message: fmt.Sprintf("baz (%s) from host: %s",
		bazResp.Version, bazResp.Host)})
	response.Write(data)
	response.Write([]byte("\n"))
	response.WriteHeader(http.StatusOK)
}
