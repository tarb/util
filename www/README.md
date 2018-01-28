# WWW

Simple Go library for making http requests

### Use

A simple library for making more complex http requests easier and neater. The library is designed to cut down on the boiler plate needed to create requests and presents a small set of chainable methods to build a request, followed up by a Collect- method to handle the response.


To change which http.Client is used for all requests (from http.DefaultClient) 
```go
SetClient( &http.Client{} )
```

Sets default values automatically used in each request
```go
SetDefaultHeaders(func(h http.Header) {
    h.Set("Header1", "Value")
    h.Set("Header2", "Value")
})
```

Perform a POST login with form data, and store the result into a struct
```go
type loginResult struct {
    Session string 
    Status  string 
}

var lResult loginResult

err := Post("https://localhost/api/login").
    WithFormBody(func(v url.Values) {
        v.Set("username", "USERNAME")
        v.Set("password", "PASSWORD")
    }).
    CollectJSON(&lResult)
```

Perform a Get request with queries, a json body and extra header values
https://google.com?q=searchterm&results=10
```go
err := Build(http.MethodGet, "https", "google.com", "").
    WithQuery(func(q url.Values) {
        q.Set("q","searchterm")
        q.Set("results","10")
    }).
    WithHeaders(func(h http.Header) {
        h.Set("Header3", "Value")
        h.Set("Header4", "Value")
    }).
    WithJSONBody(valToMarshal).
    CollectJSON(&resultStruct)
```

### TODO

* Explore adding a Builder type, to contain the reference to the default headers and client. This would allow the use of multiple instances of www