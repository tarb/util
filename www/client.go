package www

import (
	"net/http"
	"net/url"
	"time"
)

//
type Client struct {
	HTTPClient *http.Client
	MaxDelay   time.Duration
}

//
func New(cli *http.Client) *Client {
	return &Client{HTTPClient: cli}
}

//
func NewDefault() *Client {
	return &Client{HTTPClient: http.DefaultClient}
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
	u := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}

	return &httpCall{client: c.HTTPClient, req: c.newReq(method, u)}
}

// newReq builds a http.Request from the parameter values, and the
// default headers.
func (c *Client) newReq(method string, u *url.URL) *http.Request {
	return &http.Request{
		Method:        method,
		URL:           u,
		Host:          u.Host,
		Body:          nil,
		ContentLength: 0,
	}
}
