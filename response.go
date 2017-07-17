package httpfake

// Response ...
type Response struct {
	StatusCode int
	BodyBuffer []byte
}

// NewResponse ...
func NewResponse() *Response {
	return &Response{}
}

// Status ...
func (r *Response) Status(status int) *Response {
	r.StatusCode = status
	return r
}

// BodyString ...
func (r *Response) BodyString(body string) *Response {
	r.BodyBuffer = []byte(body)
	return r
}
