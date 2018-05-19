package cluster

import "net/http"

type RateLimiter struct {

}

func (l *RateLimiter) HandleRequest(writer http.ResponseWriter, request *http.Request) bool {
	return true
}
