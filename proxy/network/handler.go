package network

import "net/http"

// HTTP handler handles the upcoming requests.
type HttpHandler interface {
	HandleRequest(w http.ResponseWriter, r *http.Request) bool
}
