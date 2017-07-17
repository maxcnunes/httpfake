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
