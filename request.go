package httpfake

import (
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
)

// Request stores the settings for a request handler
// Such as how to match this handler for the incoming requests
// And how this request will respond back
type Request struct {
	sync.Mutex
	Method       string
	URL          *url.URL
	Response     *Response
	CustomHandle Responder
	assertions   []Assertor
	called       int
}

// NewRequest creates a new Request
func NewRequest() *Request {
	return &Request{
		URL:      &url.URL{},
		Response: NewResponse(),
		called:   0,
	}
}

// Get sets a GET request handler for a given path
func (r *Request) Get(path string) *Request {
	return r.method("GET", path)
}

// Post sets a POST request handler for a given path
func (r *Request) Post(path string) *Request {
	return r.method("POST", path)
}

// Put sets a PUT request handler for a given path
func (r *Request) Put(path string) *Request {
	return r.method("PUT", path)
}

// Patch sets a PATCH request handler for a given path
func (r *Request) Patch(path string) *Request {
	return r.method("PATCH", path)
}

// Delete ...
func (r *Request) Delete(path string) *Request {
	return r.method("DELETE", path)
}

// Head sets a HEAD request handler for a given path
func (r *Request) Head(path string) *Request {
	return r.method("HEAD", path)
}

// Handle sets a custom handle
// By setting this responder it gives full control to the user over this request handler
func (r *Request) Handle(handle Responder) {
	r.CustomHandle = handle
}

// Reply sets a response status for this request
// And returns the Response struct to allow chaining the response settings
func (r *Request) Reply(status int) *Response {
	return r.Response.Status(status)
}

func (r *Request) method(method, path string) *Request {
	if path != "/" {
		r.URL.Path = path
	}
	r.Method = strings.ToUpper(method)
	return r
}

func (r *Request) runAssertions(t testing.TB, testReq *http.Request) {
	for _, assertor := range r.assertions {
		assertor.Log(t)
		if err := assertor.Assert(testReq); err != nil {
			assertor.Error(t, err)
		}
	}
}

// AssertQueries will assert that the provided query parameters are present in the requests to this handler
func (r *Request) AssertQueries(key ...string) *Request {
	r.assertions = append(r.assertions, &requiredQueries{Keys: key})
	return r
}

// AssertQueryValue will assert that the provided query parameter and value are present in the requests to this handler
func (r *Request) AssertQueryValue(key, value string) *Request {
	r.assertions = append(r.assertions, &requiredQueryValue{Key: key, ExpectedValue: value})
	return r
}

// AssertHeaders will assert that the provided header keys are present in the requests to this handler
func (r *Request) AssertHeaders(keys ...string) *Request {
	r.assertions = append(r.assertions, &requiredHeaders{Keys: keys})
	return r
}

// AssertHeaderValue will assert that the provided header key and value are present in the requests to this handler
func (r *Request) AssertHeaderValue(key, value string) *Request {
	r.assertions = append(r.assertions, &requiredHeaderValue{Key: key, ExpectedValue: value})
	return r
}

// AssertBody will assert that that the provided body matches in the requests to this handler
func (r *Request) AssertBody(body []byte) *Request {
	r.assertions = append(r.assertions, &requiredBody{ExpectedBody: body})
	return r
}

// AssertCustom will run the provided assertor against requests to this handler
func (r *Request) AssertCustom(assertor Assertor) *Request {
	r.assertions = append(r.assertions, assertor)
	return r
}
