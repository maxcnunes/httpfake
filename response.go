package httpfake

import "net/http"

// Response stores the settings defined by the request handler
// of how it will respond the request back
type Response struct {
	StatusCode int
	BodyBuffer []byte
	Header     http.Header
}

// NewResponse creates a new Response
func NewResponse() *Response {
	return &Response{
		Header: make(http.Header),
	}
}

// Status sets the response status
func (r *Response) Status(status int) *Response {
	r.StatusCode = status
	return r
}

// SetHeader sets the a HTTP header to the response
func (r *Response) SetHeader(key, value string) *Response {
	r.Header.Set(key, value)
	return r
}

// AddHeader adds a HTTP header into the response
func (r *Response) AddHeader(key, value string) *Response {
	r.Header.Add(key, value)
	return r
}

// BodyString sets the response body
func (r *Response) BodyString(body string) *Response {
	r.BodyBuffer = []byte(body)
	return r
}
