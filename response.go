package httpfake

import "net/http"

// Response ...
type Response struct {
	StatusCode int
	BodyBuffer []byte
	Header     http.Header
}

// NewResponse ...
func NewResponse() *Response {
	return &Response{
		Header: make(http.Header),
	}
}

// Status ...
func (r *Response) Status(status int) *Response {
	r.StatusCode = status
	return r
}

// SetHeader ...
func (r *Response) SetHeader(key, value string) *Response {
	r.Header.Set(key, value)
	return r
}

// AddHeader ...
func (r *Response) AddHeader(key, value string) *Response {
	r.Header.Add(key, value)
	return r
}

// BodyString ...
func (r *Response) BodyString(body string) *Response {
	r.BodyBuffer = []byte(body)
	return r
}
