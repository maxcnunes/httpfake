package httpfake

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

const assertErrorTemplate = "assertion error: %s"

// Assertor provides an interface for setting assertions for http requests
type Assertor interface {
	Assert(r *http.Request) error
	Log(t testing.TB)
	Error(t testing.TB, err error)
}

// RequiredHeaders provides an Assertor for the presence of the provided http header keys
type RequiredHeaders struct {
	Keys []string
}

// Assert runs the required headers assertion against the provided request
func (h *RequiredHeaders) Assert(r *http.Request) error {
	var missingHeaders []string

	for _, key := range h.Keys {
		if value := r.Header.Get(key); len(value) == 0 {
			missingHeaders = append(missingHeaders, key)
		}
	}

	if len(missingHeaders) > 0 {
		return fmt.Errorf("missing required header(s): %s", strings.Join(missingHeaders, ", "))
	}

	return nil
}

// Log prints a testing info log for the RequiredHeaders Assertor
func (h *RequiredHeaders) Log(t testing.TB) {
	t.Log("Testing request for required headers")
}

// Error prints a testing error for the RequiredHeaders Assertor
func (h *RequiredHeaders) Error(t testing.TB, err error) {
	t.Errorf(assertErrorTemplate, err)
}

// RequiredHeaderValue provides an Assertor for a header and its expected value
type RequiredHeaderValue struct {
	Key           string
	ExpectedValue string
}

// Assert runs the required header value assertion against the provided request
func (h *RequiredHeaderValue) Assert(r *http.Request) error {
	if value := r.Header.Get(h.Key); value != h.ExpectedValue {
		return fmt.Errorf("header %s does not have the expected value; expected %s to equal %s",
			h.Key,
			value,
			h.ExpectedValue)
	}

	return nil
}

// Log prints a testing info log for the RequiredHeaderValue Assertor
func (h *RequiredHeaderValue) Log(t testing.TB) {
	t.Logf("Testing request for required header value [%s: %s]", h.Key, h.ExpectedValue)
}

// Error prints a testing error for the RequiredHeaderValue Assertor
func (h *RequiredHeaderValue) Error(t testing.TB, err error) {
	t.Errorf(assertErrorTemplate, err)
}

// RequiredQueries provides an Assertor for the presence of the provided query parameter keys
type RequiredQueries struct {
	Keys []string
}

// Assert runs the required queries assertion against the provided request
func (q *RequiredQueries) Assert(r *http.Request) error {
	queryVals := r.URL.Query()
	var missingParams []string

	for _, key := range q.Keys {
		if value := queryVals.Get(key); len(value) == 0 {
			missingParams = append(missingParams, key)
		}
	}
	if len(missingParams) > 0 {
		return fmt.Errorf("missing required query parameter(s): %s", strings.Join(missingParams, ", "))
	}

	return nil
}

// Log prints a testing info log for the RequiredQueries Assertor
func (q *RequiredQueries) Log(t testing.TB) {
	t.Log("Testing request for required query parameters")
}

// Error prints a testing error for the RequiredQueries Assertor
func (q *RequiredQueries) Error(t testing.TB, err error) {
	t.Errorf(assertErrorTemplate, err)
}

// RequiredQueryValue provides an Assertor for a query parameter and its expected value
type RequiredQueryValue struct {
	Key           string
	ExpectedValue string
}

// Assert runs the required query value assertion against the provided request
func (q *RequiredQueryValue) Assert(r *http.Request) error {
	if value := r.URL.Query().Get(q.Key); value != q.ExpectedValue {
		return fmt.Errorf("query %s does not have the expected value; expected %s to equal %s", q.Key, value, q.ExpectedValue)
	}
	return nil
}

// Log prints a testing info log for the RequiredQueryValue Assertor
func (q *RequiredQueryValue) Log(t testing.TB) {
	t.Logf("Testing request for required query parameter value [%s: %s]", q.Key, q.ExpectedValue)
}

// Error prints a testing error for the RequiredQueryValue Assertor
func (q *RequiredQueryValue) Error(t testing.TB, err error) {
	t.Errorf(assertErrorTemplate, err)
}

// RequiredBody provides an Assertor for the expected value of the request body
type RequiredBody struct {
	ExpectedBody []byte
}

// Assert runs the required body assertion against the provided request
func (b *RequiredBody) Assert(r *http.Request) error {
	if r.Body == nil {
		return fmt.Errorf("error reading the request body; the request body is nil")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading the request body: %w", err)
	}

	if !bytes.EqualFold(b.ExpectedBody, body) {
		return fmt.Errorf("request body does not have the expected value; expected %s to equal %s",
			string(body[:]),
			string(b.ExpectedBody[:]))
	}

	return nil
}

// Log prints a testing info log for the RequiredBody Assertor
func (b *RequiredBody) Log(t testing.TB) {
	t.Log("Testing request for required a required body")
}

// Error prints a testing error for the RequiredBody Assertor
func (b *RequiredBody) Error(t testing.TB, err error) {
	t.Errorf(assertErrorTemplate, err)
}
