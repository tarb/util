package www

import (
	"net/http"
	"sync"
)

var dhMut sync.RWMutex
var defaultHeaders = make(http.Header)

// SetDefaultHeaders gives an interface for setting which headers
// will be added by defualt when doing a new httpCall.
//
// Example
// SetDefaultHeaders(func(h http.Header) {
//     h.Set("User-Agent", "CustomAgentString")
// })
func SetDefaultHeaders(fn func(http.Header)) {
	dhMut.Lock()
	fn(defaultHeaders)
	dhMut.Unlock()
}

// ClearDefaultHeaders clears any values stored inside the default
// headers.
func ClearDefaultHeaders() {
	dhMut.Lock()
	defaultHeaders = make(http.Header)
	dhMut.Unlock()
}

// getDefaultHeaders returns a new http.Header with a copy of the
// values that have been set as default.
func getDefaultHeaders() http.Header {
	var h = make(http.Header)

	dhMut.RLock()
	for k, v := range defaultHeaders {
		h[k] = v
	}
	dhMut.RUnlock()

	return h
}
