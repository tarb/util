package www

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

//
type Client struct {
	mu         sync.RWMutex
	headers    http.Header // default headers
	HTTPClient *http.Client

	MaxDelay time.Duration
}

//
func New(cli *http.Client) *Client {
	return &Client{
		headers:    make(http.Header),
		HTTPClient: cli,
	}
}

//
func NewDefault() *Client {
	return &Client{
		headers:    make(http.Header),
		HTTPClient: http.DefaultClient,
	}
}

// Get creates a new httpCall with default values and
// using the GET http method
func (c *Client) Get(urlStr string) *httpCall {
	u, err := url.Parse(urlStr)

	return &httpCall{client: c.HTTPClient, err: err, req: c.newReq(http.MethodGet, u)}
}

// Post creates a new httpCall with default values and
// using the POST http method
func (c *Client) Post(urlStr string) *httpCall {
	u, err := url.Parse(urlStr)

	return &httpCall{client: c.HTTPClient, err: err, req: c.newReq(http.MethodPost, u)}
}

// Build creates a custom httpCall from the given parameters;
// scheme, host and path are used to create a new url.URL value.
//
// There is no validation that the method value is a correct http
// method, suggest using contants found in http package
// (http.MethodGet, http.MethodPost, etc)
func (c *Client) Build(method, scheme, host, path string) *httpCall {
	var u = &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}

	return &httpCall{client: c.HTTPClient, req: c.newReq(method, u)}
}

// SetDefaultHeaders gives an interface for setting which headers
// will be added by defualt when doing a new httpCall.
//
// Example
// SetDefaultHeaders(func(h http.Header) {
//     h.Set("User-Agent", "CustomAgentString")
// })
func (c *Client) SetDefaultHeaders(fn func(http.Header)) {
	c.mu.Lock()
	fn(c.headers)
	c.mu.Unlock()
}

// ClearDefaultHeaders clears any values stored inside the default
// headers.
func (c *Client) ClearDefaultHeaders() {
	c.mu.Lock()
	c.headers = make(http.Header)
	c.mu.Unlock()
}

// getDefaultHeaders returns a new http.Header with a copy of the
// values that have been set as default.
func (c *Client) getDefaultHeaders() http.Header {
	var h = make(http.Header)

	c.mu.RLock()
	for k, v := range c.headers {
		h[k] = v
	}
	c.mu.RUnlock()

	return h
}

// newReq builds a http.Request from the parameter values, and the
// default headers.
func (c *Client) newReq(method string, u *url.URL) *http.Request {
	return &http.Request{
		Method:        method,
		URL:           u,
		Host:          u.Host,
		Header:        c.getDefaultHeaders(),
		Body:          nil,
		ContentLength: 0,
	}
}
