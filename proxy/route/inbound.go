package route

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/sczyh30/waffle-mesh/proxy/network/core"
)

type InboundRouter struct {

}

func NewInboundRouter() *InboundRouter {
	return &InboundRouter{}
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

	client := core.NewHttpClient()
	targetUrl := "http://localhost:" + fmt.Sprint(port) + request.RequestURI
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
