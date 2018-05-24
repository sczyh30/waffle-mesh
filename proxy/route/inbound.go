package route

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/sczyh30/waffle-mesh/proxy/network/core"
)

type InboundRouter struct {
	inboundPort uint32
}

func NewInboundRouter(port uint32) *InboundRouter {
	return &InboundRouter{inboundPort: port}
}

func (r *InboundRouter) HandleRequest(writer http.ResponseWriter, request *http.Request) bool {
	// Parse the target port.
	h := strings.Split(request.Host, ":")
	var port int
	if len(h) == 1 {
		port = 80
	} else {
		port, _ = strconv.Atoi(h[1])
	}

	if port == int(r.inboundPort) {
		writer.WriteHeader(503)
		fmt.Fprintf(writer, "Invalid target port: %d\n", r.inboundPort)
		return StopChain
	}

	client := core.NewHttpClient()
	targetUrl := "http://localhost:" + fmt.Sprint(port) + request.RequestURI
	log.Printf("[Inbound Server] Will redirect to: %s\n", targetUrl)
	newRequest, _ := http.NewRequest(request.Method, targetUrl, request.Body)

	newRequest.Header = request.Header
	resp, err := client.Do(newRequest)
	if err != nil {
		writer.WriteHeader(503)
		fmt.Fprint(writer, "service unavailable: " + err.Error())
		return StopChain
	}

	writer.WriteHeader(resp.StatusCode)
	io.Copy(writer, resp.Body)

	return StopChain
}
