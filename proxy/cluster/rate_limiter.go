package cluster

import "net/http"

type RateLimitHandler struct {

}

func (l *RateLimitHandler) HandleRequest(writer http.ResponseWriter, request *http.Request) bool {

	return true
}
