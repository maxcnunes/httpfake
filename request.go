package httpfake

import (
	"net/url"
	"strings"
)

// Request stores the settings for a request handler
// Such as how to match this handler for the incoming requests
// And how this request will respond back
type Request struct {
	Method       string
	URL          *url.URL
	Response     *Response
	CustomHandle Responder
}

// NewRequest creates a new Request
func NewRequest() *Request {
	return &Request{
		URL:      &url.URL{},
		Response: NewResponse(),
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
// And returns the Reponse struct to allow chaining the response settings
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
