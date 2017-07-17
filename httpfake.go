package httpfake

import (
	"net/http"
	"net/http/httptest"
)

// HTTPFake ...
type HTTPFake struct {
	Server          *httptest.Server
	RequestHandlers []*Request
}

// New ...
func New() *HTTPFake {
	fake := &HTTPFake{
		RequestHandlers: []*Request{},
	}

	fake.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rh := fake.findHandler(r)
		if rh == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		respond(rh, w)
	}))

	return fake
}

// NewHandler ...
func (f *HTTPFake) NewHandler() *Request {
	rh := NewRequest()
	f.RequestHandlers = append(f.RequestHandlers, rh)
	return rh
}

// Reset ...
func (f *HTTPFake) Reset() *HTTPFake {
	f.RequestHandlers = []*Request{}
	return f
}

func (f *HTTPFake) findHandler(r *http.Request) *Request {
	founds := []*Request{}
	for _, rh := range f.RequestHandlers {
		if rh.URL.String() == r.URL.String() {
			return rh
		}
		// fallback if the income request has query strings
		// and there is handlers only for the path
		if rh.URL.EscapedPath() == r.URL.EscapedPath() {
			founds = append(founds, rh)
		}
	}
	// only use the fallback if could find only one match
	if len(founds) == 1 {
		return founds[0]
	}
	return nil
}

func respond(rh *Request, w http.ResponseWriter) {
	if rh.Response.StatusCode > 0 {
		w.WriteHeader(rh.Response.StatusCode)
	}
	if len(rh.Response.BodyBuffer) > 0 {
		w.Write(rh.Response.BodyBuffer) // nolint
	}
}
