package httpfake

import (
	"net/url"
	"strings"
)

// Request ...
type Request struct {
	Method   string
	URL      *url.URL
	Response *Response
}

// NewRequest ...
func NewRequest() *Request {
	return &Request{
		URL:      &url.URL{},
		Response: NewResponse(),
	}
}

// Get ...
func (r *Request) Get(path string) *Request {
	return r.method("GET", path)
}

// Post ...
func (r *Request) Post(path string) *Request {
	return r.method("POST", path)
}

// Put ...
func (r *Request) Put(path string) *Request {
	return r.method("PUT", path)
}

// Patch ...
func (r *Request) Patch(path string) *Request {
	return r.method("PATCH", path)
}

// Delete ...
func (r *Request) Delete(path string) *Request {
	return r.method("DELETE", path)
}

// Head ...
func (r *Request) Head(path string) *Request {
	return r.method("HEAD", path)
}

// Reply ...
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
