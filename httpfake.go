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
		if rh.CustomHandle != nil {
			rh.CustomHandle(w, r, rh)
			return
		}
		Respond(w, r, rh)
	}))

	return fake
}

// NewHandler ...
func (f *HTTPFake) NewHandler() *Request {
	rh := NewRequest()
	f.RequestHandlers = append(f.RequestHandlers, rh)
	return rh
}

// ResolveURL resolves the full URL to the fake server for a given path
func (f *HTTPFake) ResolveURL(path string) string {
	return f.Server.URL + path
}

// Reset ...
func (f *HTTPFake) Reset() *HTTPFake {
	f.RequestHandlers = []*Request{}
	return f
}

func (f *HTTPFake) findHandler(r *http.Request) *Request {
	founds := []*Request{}
	url := r.URL.String()
	for _, rh := range f.RequestHandlers {
		if rh.Method != r.Method {
			continue
		}

		if rh.URL.String() == url {
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
