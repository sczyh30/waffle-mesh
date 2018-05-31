package cluster

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/time/rate"
)

const defaultLimitQps = 1000

var RateLimitConfig map[string]int

type RateLimitHandler struct {

}

func (l *RateLimitHandler) HandleRequest(writer http.ResponseWriter, request *http.Request) bool {
	arr := strings.Split(request.Host, ":")
	if len(arr) > 1 {
		clusterName := arr[0]
		if e, exists := runtimeClusterRegistry.registryMap[clusterName]; exists {
			if e.rateLimiter == nil {
				e.rateLimiter = createLimiterFor(clusterName)
			}
			if e.rateLimiter.Allow() {
				return true
			} else {
				l.handleError(writer, errors.New("cannot allow more requests (rate-limit)"), 503)
				return false
			}
		} else {
			l.handleError(writer, errors.New("cannot find the cluster for rate-limit"), 404)
			return false
		}
	}
	return true
}

func (l *RateLimitHandler) handleError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	fmt.Fprint(w, err.Error())
}

func createLimiterFor(cluster string) *rate.Limiter {
	if limit, exists := RateLimitConfig[cluster]; exists {
		return rate.NewLimiter(rate.Limit(limit), 0)
	} else {
		return rate.NewLimiter(defaultLimitQps, 0)
	}
}
