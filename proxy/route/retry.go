package route

import (
	"math"
	"net/http"
	"strings"

	"github.com/sczyh30/waffle-mesh/api/gen"
)

const (
	ServerFail = "server_fail"
	ServerConnectTimeout = "server_connect_timeout"
)

type RetryPolicy interface {
	AskForRetry() bool

	NextRetryTimeout() float64

	GetAttempts() int
}

func NewRetryPolicy(s *api.RouteAction_RetryStrategy) RetryPolicy {
	if s == nil {
		return nil
	}
	return NewExponentialBackoffRetryPolicy(int(s.RetryTimes), s.RetryTimeout.Seconds * 1000)
}

func NewFixedBackoffRetryPolicy(max int, timeout int64) RetryPolicy {
	// TODO: validate the parameters.
	return &FixedBackoffRetryPolicy{
		max: max,
		attempts: 0,
		timeout: timeout,
	}
}

type FixedBackoffRetryPolicy struct {
	max int
	attempts int

	timeout int64
}

func (s *FixedBackoffRetryPolicy) GetAttempts() int {
	return s.attempts
}

func (s *FixedBackoffRetryPolicy) AskForRetry() bool {
	if s.attempts < s.max {
		s.attempts++
		return true
	}
	return false
}

func (s *FixedBackoffRetryPolicy) NextRetryTimeout() float64 {
	return float64(s.timeout)
}

func NewExponentialBackoffRetryPolicy(max int, timeout int64) RetryPolicy {
	// TODO: validate the parameters.
	return &ExponentialBackoffRetryPolicy{
		max: max,
		attempts: 0,
		baseTimeout: timeout,
	}
}

type ExponentialBackoffRetryPolicy struct {
	max int
	attempts int

	baseTimeout int64
}

func (s *ExponentialBackoffRetryPolicy) GetAttempts() int {
	return s.attempts
}

func (s *ExponentialBackoffRetryPolicy) AskForRetry() bool {
	if s.attempts < s.max {
		s.attempts++
		return true
	}
	return false
}

func (s *ExponentialBackoffRetryPolicy) NextRetryTimeout() float64 {
	return float64(s.baseTimeout) * (math.Pow(2, float64(s.attempts)) - 1)
}

func shouldRetry(s *api.RouteAction_RetryStrategy, resp *http.Response, err error) bool {
	if s == nil {
		return false
	}
	switch s.RetryType {
	case ServerFail:
		return isTimeout(err) || isServer5xxFail(resp)
	case ServerConnectTimeout:
		return isTimeout(err)
	}
	return false
}

func isTimeout(err error) bool {
	return err != nil && strings.Contains(err.Error(), "connection timed out")
}

func isServer5xxFail(resp *http.Response) bool {
	status := resp.StatusCode
	return status >= http.StatusInternalServerError && status <= http.StatusGatewayTimeout
}
