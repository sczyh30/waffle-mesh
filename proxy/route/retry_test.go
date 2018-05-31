package route

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/duration"

	"github.com/sczyh30/waffle-mesh/api/gen"
)

type MockResponse struct {
	resp *http.Response
	err error
}

type MockClient struct {
	times int
	cur int
	respList []MockResponse
}

func newMockClient(times int) *MockClient {
	return &MockClient{
		times: times,
		cur: 0,
		respList: make([]MockResponse, times),
	}
}

func (c *MockClient) ReturnAt(t int, resp *http.Response, err error) {
	if t > c.times - 1 {
		return
	}
	c.respList[t] = MockResponse{resp: resp, err: err}
}

func (c *MockClient) DoRequest() (*http.Response, error) {
	resp := c.respList[c.cur]
	c.cur++
	return resp.resp, resp.err
}

func TestServerFailRetry(t *testing.T) {
	client := newMockClient(3)
	client.ReturnAt(0, nil, errors.New("error: connection timed out"))
	client.ReturnAt(1, &http.Response{StatusCode: 500}, nil)
	client.ReturnAt(2, &http.Response{StatusCode: 200}, nil)

	strategy := &api.RouteAction_RetryStrategy{RetryType:ServerFail, RetryTimes: 2, RetryTimeout: &duration.Duration{Seconds: 1}}
	retryPolicy := NewRetryPolicy(strategy)
	response, err := client.DoRequest()
	for {
		if shouldRetry(strategy, response, err) {
			if retryPolicy.AskForRetry() {
				retryTimeout := retryPolicy.NextRetryTimeout()
				fmt.Printf("This time retry will wait for %f ms\n", retryTimeout)
				// Sleep for timeouts.
				time.Sleep(time.Millisecond * time.Duration(retryTimeout))

				// Do retry.
				response, err = client.DoRequest()
			} else {
				// Retry times exceeded.
				t.Errorf("service unavailable: retry exceeded (%d)\n", retryPolicy.GetAttempts())
				return
			}
		} else {
			break
		}
	}
}
