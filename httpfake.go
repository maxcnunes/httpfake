// Package httpfake provides a simple wrapper for httptest
// with a handful chainable API for setting up handlers to a fake server.
// This package is aimed to be used in tests where the original external server
// must not be reached. Instead is used in its place a fake server
// which can be configured to handle any request as desired.
package httpfake

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	netURL "net/url"
	"strings"
)

// HTTPFake is the root struct for the fake server
type HTTPFake struct {
	Server          *httptest.Server
	RequestHandlers []*Request
}

// New starts a httptest.Server as the fake server
// and sets up the initial configuration to this server's request handlers
func New() *HTTPFake {
	fake := &HTTPFake{
		RequestHandlers: []*Request{},
	}

	fake.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rh, err := fake.findHandler(r)
		if err != nil {
			printError(fmt.Sprintf("error finding handler: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if rh == nil {
			errMsg := fmt.Sprintf(
				"not found request handler for [%s: %s]; registered handlers are:\n",
				r.Method, r.URL,
			)
			for _, frh := range fake.RequestHandlers {
				errMsg += fmt.Sprintf("* [%s: %s]\n", frh.Method, frh.URL.Path)
			}
			printError(errMsg)
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

// NewHandler initializes the configuration for a new request handler
func (f *HTTPFake) NewHandler() *Request {
	rh := NewRequest()
	f.RequestHandlers = append(f.RequestHandlers, rh)
	return rh
}

// ResolveURL resolves the full URL to the fake server for a given path
func (f *HTTPFake) ResolveURL(path string, args ...interface{}) string {
	format := f.Server.URL + path
	return fmt.Sprintf(format, args...)
}

// Reset wipes the request handlers definitions
func (f *HTTPFake) Reset() *HTTPFake {
	f.RequestHandlers = []*Request{}
	return f
}

func (f *HTTPFake) findHandler(r *http.Request) (*Request, error) {
	founds := []*Request{}
	url := r.URL.String()
	path := getURLPath(url)
	for _, rh := range f.RequestHandlers {
		if rh.Method != r.Method {
			continue
		}

		rhURL, err := netURL.QueryUnescape(rh.URL.String())
		if err != nil {
			return nil, err
		}

		if rhURL == url {
			return rh, nil
		}

		// fallback if the income request has query strings
		// and there is handlers only for the path
		if getURLPath(rhURL) == path {
			founds = append(founds, rh)
		}
	}
	// only use the fallback if could find only one match
	if len(founds) == 1 {
		return founds[0], nil
	}

	return nil, nil
}

func getURLPath(url string) string {
	return strings.Split(url, "?")[0]
}

func printError(msg string) {
	fmt.Println("\033[0;31mhttpfake: " + msg + "\033[0m")
}
