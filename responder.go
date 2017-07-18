package httpfake

import "net/http"

// Responder are callbacks to handle the request and write the response
type Responder func(w http.ResponseWriter, r *http.Request, rh *Request)

// Respond writes the response based in the request handler settings
func Respond(w http.ResponseWriter, r *http.Request, rh *Request) {
	if len(rh.Response.Header) > 0 {
		for k := range rh.Response.Header {
			w.Header().Add(k, rh.Response.Header.Get(k))
		}
	}
	if rh.Response.StatusCode > 0 {
		w.WriteHeader(rh.Response.StatusCode)
	}
	if len(rh.Response.BodyBuffer) > 0 {
		w.Write(rh.Response.BodyBuffer) // nolint
	}
}
