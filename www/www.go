package www

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// thrown errors
var (
	ErrMaxAttempts = errors.New("http request failed: max attempts reached")
)

// httpCall contains the building of the requestm abd the
// resulting response. Its not exported so that users use
// it in the intended way - culminating each request in a
// CollectX call.
type httpCall struct {
	err     error
	body    []byte
	client  *http.Client
	req     *http.Request
	resp    *http.Response
	headers http.Header
}

// WithQuery is used to update the url query it expects a func
// with a url.Values arg, allowing you to use Values methods to
// build the url query.
//
// Example
// Get("https://localhost/search").
// WithQuery(func(q url.Values) {
//     q.Set("q","golang")
//     q.Set("results","10")
// }). ...
//
// will produce:  https://localhost/search?q=golang&results=10
func (c *httpCall) WithQuery(fn func(url.Values)) *httpCall {
	queries := make(url.Values)
	fn(queries)

	c.req.URL.RawQuery = queries.Encode()
	return c
}

// WithJSONBody Marshals the value into json, and updates
// the interal requests Body to contain a reader with the
// serialized bytes, the requests Contentlength is also
// updated to reflect the length of the serialized bytes,
//
// This method also adds the application/json Content-Type
// header to the request
func (c *httpCall) WithJSONBody(jsonBody interface{}) *httpCall {
	c.body, c.err = json.Marshal(jsonBody)

	c.req.Body = ioutil.NopCloser(bytes.NewReader(c.body))
	c.req.ContentLength = int64(len(c.body))
	if c.req.Header == nil {
		c.req.Header = make(http.Header)
	}
	c.req.Header.Set("Content-Type", "application/json; charset=utf-8")

	return c
}

// WithFormBody is used to set the body of a request to contain
// the formatted form values. it expects a func with a url.Values
// arg, allowing you to use Values methods to build the form data.
//
// Example
// Post("https://localhost/login").
// WithFormBody(func(form url.Values) {
//     form.Set("username","admin")
//     form.Set("password","hunter42")
// }).
//
// This method also adds the application/x-www-form-urlencoded
// Content-Type header to the request
func (c *httpCall) WithFormBody(fn func(url.Values)) *httpCall {
	formData := make(url.Values)
	fn(formData)

	c.body = []byte(formData.Encode())
	c.req.Body = ioutil.NopCloser(bytes.NewReader(c.body))
	c.req.ContentLength = int64(len(c.body))

	if c.req.Header == nil {
		c.req.Header = make(http.Header)
	}
	c.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c
}

// WithTextBody sets the requests body to be a reader drawing from
// the string parameter
//
// This method also adds the text/plain Content-Type header to the
// request
func (c *httpCall) WithTextBody(body string) *httpCall {
	c.body = []byte(body)
	c.req.Body = ioutil.NopCloser(bytes.NewReader(c.body))
	c.req.ContentLength = int64(len(c.body))

	if c.req.Header == nil {
		c.req.Header = make(http.Header)
	}
	c.req.Header.Set("Content-Type", "text/plain")

	return c
}

// WithHeaders gives an interface for adding new headers to the
// request. The header value will already contain any set default
// default values, and may contain headers set from other methods,
// primarily Content-Type headers set when adding a body to the req
//
// Example
// Get("https://google.com").
//     WithHeaders(func(h http.Header) {
//         h.Set("User-Agent", "CustomAgentString")
//     }).
func (c *httpCall) WithHeaders(fn func(http.Header)) *httpCall {
	if c.req.Header == nil {
		c.req.Header = make(http.Header)
	}

	fn(c.req.Header)

	return c
}

// CollectJSON finalizes the httpCall and Unmarshals the responses
// body into the parameter value. If somehwere in the httpCall method
// chain has returned an error, dont run the request and return the
// error. If the request returns an error, return that error - or if
// the responses status code is not 2--, return a new error detailing
// the responses status
func (c *httpCall) CollectJSON(obj interface{}) error {
	defer func() {
		if c.resp != nil && c.resp.Body != nil {
			c.resp.Body.Close()
		}
	}()

	if c.err != nil {
		return c.err
	}

	c.err = json.NewDecoder(c.resp.Body).Decode(obj)

	return c.err
}

// CollectXML finalizes the httpCall and Unmarshals the responses
// body into the parameter value. If somewhere in the httpCall method
// chain has returned an error, dont run the request and return the
// error. If the request returns an error, return that error - or if
// the responses status code is not 2--, return a new error detailing
// the responses status
func (c *httpCall) CollectXML(obj interface{}) error {
	defer func() {
		if c.resp != nil && c.resp.Body != nil {
			c.resp.Body.Close()
		}
	}()

	if c.err != nil {
		return c.err
	}

	c.err = xml.NewDecoder(c.resp.Body).Decode(obj)

	return c.err
}

// CollectString finalizes the httpCall and returns the response body
// as a string. If somehwere in the httpCall method
// chain has returned an error, dont run the request and return the
// error. If the request returns an error, return that error - or if
// the responses status code is not 2--, return a new error detailing
// the responses status
func (c *httpCall) CollectString() (string, error) {
	defer func() {
		if c.resp != nil && c.resp.Body != nil {
			c.resp.Body.Close()
		}
	}()

	if c.err != nil {
		return "", c.err
	}

	var bs []byte
	bs, c.err = ioutil.ReadAll(c.resp.Body)

	return string(bs), c.err
}

// CollectBytes finalizes the httpCall and returns the response body
// as an array of bytes. If somehwere in the httpCall method
// chain has returned an error, dont run the request and return the
// error. If the request returns an error, return that error - or if
// the responses status code is not 2--, return a new error detailing
// the responses status
func (c *httpCall) CollectBytes() ([]byte, error) {
	defer func() {
		if c.resp != nil && c.resp.Body != nil {
			c.resp.Body.Close()
		}
	}()

	if c.err != nil {
		return nil, c.err
	}

	var bs []byte
	bs, c.err = ioutil.ReadAll(c.resp.Body)

	return bs, c.err
}

// CollectResponse finalizes the httpCall. If somewhere in the httpCall
// method chain has returned an error, dont run the request and return
// the  error. If the request returns an error, return that error - or
// if the responses status code is not 2--, return a new error detailing
// the responses status
func (c *httpCall) CollectResponse() (*http.Response, error) {
	if c.err != nil {
		return nil, c.err
	}

	return c.resp, nil
}

// do checks for errors in the construction of the req, and if not
// present runs the request, storing the value and error into itself
func (c *httpCall) Do() *httpCall {
	// if err has already been set by something
	// dont bother trying to do the request
	if c.err != nil {
		return c
	}

	// fmt.Println(c.req.URL.String())

	c.resp, c.err = c.client.Do(c.req)

	if c.err == nil {
		// check the encoding, and wrap in gzip if needed
		if !strings.Contains(c.resp.Header.Get("Accept-Encoding"), "gzip") {
			c.resp.Body, c.err = gzip.NewReader(c.resp.Body)
		}
		// client errors should not change just with retrying
		// so cancel out straight away
		if c.resp.StatusCode/100 == 4 {
			bs, _ := ioutil.ReadAll(c.resp.Body)
			c.err = StatusError{Status: c.resp.Status, StatusCode: c.resp.StatusCode, Body: bs}
		}
	}

	return c
}

//
func (c *httpCall) DoWithRetry(maxAttempts int, delay DelayFunc) *httpCall {
	// if err has already been set by something
	// dont bother trying to do the request
	if c.err != nil {
		return c
	}

	attempts := 0

	for (attempts == 0 || c.err != nil) && attempts < maxAttempts {

		if attempts > 0 {
			time.Sleep(delay(attempts))
		}

		// set the request body to read from the saved bytes
		if c.body != nil {
			c.req.Body = ioutil.NopCloser(bytes.NewReader(c.body))
		}

		c.resp, c.err = c.client.Do(c.req)

		if c.err == nil {
			// client errors should not change just with retrying, so cancel out
			// straight away
			if code := c.resp.StatusCode; code/100 == 4 {
				bs, _ := ioutil.ReadAll(c.resp.Body)
				c.err = StatusError{Status: c.resp.Status, StatusCode: c.resp.StatusCode, Body: bs}
			}

			return c
		}

		attempts++
	}

	if attempts >= maxAttempts {
		c.err = ErrMaxAttempts
	}

	return c
}
