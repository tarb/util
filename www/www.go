package www

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

// httpCall contains the building of the requestm abd the
// resulting response. Its not exported so that users use
// it in the intended way - culminating each request in a
// CollectX call.
type httpCall struct {
	err    error
	client *http.Client
	req    *http.Request
	resp   *http.Response
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
	var queries = make(url.Values)
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
	var bs []byte
	bs, c.err = json.Marshal(jsonBody)

	c.req.Body = ioutil.NopCloser(bytes.NewReader(bs))
	c.req.ContentLength = int64(len(bs))
	c.req.Header.Set("Content-Type", "application/json")
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
	var formData = make(url.Values)
	fn(formData)

	var bs = []byte(formData.Encode())

	c.req.Body = ioutil.NopCloser(bytes.NewReader(bs))
	c.req.ContentLength = int64(len(bs))
	c.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c
}

// WithTextBody sets the requests body to be a reader drawing from
// the string parameter
//
// This method also adds the text/plain Content-Type header to the
// request
func (c *httpCall) WithTextBody(body string) *httpCall {
	var bs = []byte(body)

	c.req.Body = ioutil.NopCloser(bytes.NewReader(bs))
	c.req.ContentLength = int64(len(bs))
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
	if c.do(); c.err != nil {
		return c.err
	}

	if c.resp.StatusCode/100 == 2 {
		c.err = json.NewDecoder(c.resp.Body).Decode(obj)
	} else {
		var bs []byte
		bs, c.err = ioutil.ReadAll(c.resp.Body)
		c.err = StatusError{Status: c.resp.Status, StatusCode: c.resp.StatusCode, Body: string(bs)}
	}

	c.resp.Body.Close()
	return c.err
}

// CollectXML finalizes the httpCall and Unmarshals the responses
// body into the parameter value. If somewhere in the httpCall method
// chain has returned an error, dont run the request and return the
// error. If the request returns an error, return that error - or if
// the responses status code is not 2--, return a new error detailing
// the responses status
func (c *httpCall) CollectXML(obj interface{}) error {
	if c.do(); c.err != nil {
		return c.err
	}

	if c.resp.StatusCode/100 == 2 {
		c.err = xml.NewDecoder(c.resp.Body).Decode(obj)
	} else {
		var bs []byte
		bs, c.err = ioutil.ReadAll(c.resp.Body)
		c.err = StatusError{Status: c.resp.Status, StatusCode: c.resp.StatusCode, Body: string(bs)}
	}

	c.resp.Body.Close()
	return c.err
}

// CollectJSON finalizes the httpCall and returns the response body
// as a string. If somehwere in the httpCall method
// chain has returned an error, dont run the request and return the
// error. If the request returns an error, return that error - or if
// the responses status code is not 2--, return a new error detailing
// the responses status
func (c *httpCall) CollectString() (string, error) {
	if c.do(); c.err != nil {
		return "", c.err
	}

	var bs []byte
	if c.resp.StatusCode/100 == 2 {
		bs, c.err = ioutil.ReadAll(c.resp.Body)
	} else {
		bs, c.err = ioutil.ReadAll(c.resp.Body)
		c.err = StatusError{Status: c.resp.Status, StatusCode: c.resp.StatusCode, Body: string(bs)}
	}

	c.resp.Body.Close()

	return string(bs), c.err
}

// CollectResponse finalizes the httpCall. If somewhere in the httpCall
// method chain has returned an error, dont run the request and return
// the  error. If the request returns an error, return that error - or
// if the responses status code is not 2--, return a new error detailing
// the responses status
func (c *httpCall) CollectResponse() (*http.Response, error) {
	if c.do(); c.err != nil {
		return nil, c.err
	}

	if c.resp.StatusCode/200 != 2 {
		c.err = StatusError{Status: c.resp.Status, StatusCode: c.resp.StatusCode}
	}

	return c.resp, nil
}

// do checks for errors in the construction of the req, and if not
// present runs the request, storing the value and error into itself
func (c *httpCall) do() {
	// if err has already been set by something
	// dont bother trying to do the request
	if c.err != nil {
		return
	}

	c.resp, c.err = c.client.Do(c.req)
}
