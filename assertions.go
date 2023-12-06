package httpfake

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/tidwall/gjson"
)

const assertErrorTemplate = "assertion error: %s"

// Assertor provides an interface for setting assertions for http requests.
type Assertor interface {
	Assert(r *http.Request) error
	Log(t testing.TB)
	Error(t testing.TB, err error)
}

// requiredHeaders provides an Assertor for the presence of the provided http header keys.
type requiredHeaders struct {
	Keys []string
}

// Assert runs the required headers assertion against the provided request.
func (h *requiredHeaders) Assert(r *http.Request) error {
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

// Log prints a testing info log for the requiredHeaders Assertor.
func (h *requiredHeaders) Log(t testing.TB) {
	t.Log("Testing request for required headers")
}

// Error prints a testing error for the requiredHeaders Assertor.
func (h *requiredHeaders) Error(t testing.TB, err error) {
	t.Errorf(assertErrorTemplate, err)
}

// requiredHeaderValue provides an Assertor for a header and its expected value.
type requiredHeaderValue struct {
	Key           string
	ExpectedValue string
}

// Assert runs the required header value assertion against the provided request.
func (h *requiredHeaderValue) Assert(r *http.Request) error {
	if value := r.Header.Get(h.Key); value != h.ExpectedValue {
		return fmt.Errorf("header %s does not have the expected value; expected %s to equal %s",
			h.Key,
			value,
			h.ExpectedValue)
	}

	return nil
}

// Log prints a testing info log for the requiredHeaderValue Assertor.
func (h *requiredHeaderValue) Log(t testing.TB) {
	t.Logf("Testing request for a required header value [%s: %s]", h.Key, h.ExpectedValue)
}

// Error prints a testing error for the requiredHeaderValue Assertor.
func (h *requiredHeaderValue) Error(t testing.TB, err error) {
	t.Errorf(assertErrorTemplate, err)
}

// requiredQueries provides an Assertor for the presence of the provided query parameter keys.
type requiredQueries struct {
	Keys []string
}

// Assert runs the required queries assertion against the provided request.
func (q *requiredQueries) Assert(r *http.Request) error {
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

// Log prints a testing info log for the requiredQueries Assertor.
func (q *requiredQueries) Log(t testing.TB) {
	t.Log("Testing request for required query parameters")
}

// Error prints a testing error for the requiredQueries Assertor.
func (q *requiredQueries) Error(t testing.TB, err error) {
	t.Errorf(assertErrorTemplate, err)
}

// requiredQueryValue provides an Assertor for a query parameter and its expected value.
type requiredQueryValue struct {
	Key           string
	ExpectedValue string
}

// Assert runs the required query value assertion against the provided request.
func (q *requiredQueryValue) Assert(r *http.Request) error {
	if value := r.URL.Query().Get(q.Key); value != q.ExpectedValue {
		return fmt.Errorf("query %s does not have the expected value; expected %s to equal %s", q.Key, value, q.ExpectedValue)
	}
	return nil
}

// Log prints a testing info log for the requiredQueryValue Assertor.
func (q *requiredQueryValue) Log(t testing.TB) {
	t.Logf("Testing request for a required query parameter value [%s: %s]", q.Key, q.ExpectedValue)
}

// Error prints a testing error for the requiredQueryValue Assertor.
func (q *requiredQueryValue) Error(t testing.TB, err error) {
	t.Errorf(assertErrorTemplate, err)
}

// requiredBody provides an Assertor for the expected value of the request body.
type requiredBody struct {
	ExpectedBody []byte
}

// Assert runs the required body assertion against the provided request.
func (b *requiredBody) Assert(r *http.Request) error {
	if r.Body == nil {
		return fmt.Errorf("error reading the request body; the request body is nil")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading the request body: %s", err.Error())
	}

	if !bytes.EqualFold(b.ExpectedBody, body) {
		return fmt.Errorf("request body does not have the expected value; expected %s to equal %s",
			string(body),
			string(b.ExpectedBody))
	}

	return nil
}

// Log prints a testing info log for the requiredBody Assertor.
func (b *requiredBody) Log(t testing.TB) {
	t.Log("Testing request for a required body value")
}

// Error prints a testing error for the requiredBody Assertor.
func (b *requiredBody) Error(t testing.TB, err error) {
	t.Errorf(assertErrorTemplate, err)
}

// CustomAssertor provides a function signature that implements the Assertor interface. This allows for
// adhoc creation of a custom assertion for use with the AssertCustom assertor.
type CustomAssertor func(r *http.Request) error

// Assert runs the CustomAssertor assertion against the provided request.
func (c CustomAssertor) Assert(r *http.Request) error {
	return c(r)
}

// Log prints a testing info log for the CustomAssertor.
func (c CustomAssertor) Log(t testing.TB) {
	t.Log("Testing request with a custom assertor")
}

// Error prints a testing error for the CustomAssertor.
func (c CustomAssertor) Error(t testing.TB, err error) {
	t.Errorf(assertErrorTemplate, err)
}

type subJSON struct {
	necessaryFields map[string]any
}

type gjsonForEachFunc func(key, value gjson.Result) bool

func (j *subJSON) fillFunc(keyPrefix string) gjsonForEachFunc {
	return func(key, value gjson.Result) bool {
		keyString := key.String()

		if keyPrefix != "" {
			keyString = fmt.Sprintf("%s.%s", keyPrefix, keyString)
		}

		if value.IsArray() || value.IsObject() {
			value.ForEach(j.fillFunc(keyString))
			return true
		}

		j.necessaryFields[keyString] = value.Value()

		return true
	}
}

func (j *subJSON) Assert(r *http.Request) error {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	_ = r.Body.Close()

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	js := gjson.ParseBytes(bodyBytes)

	var assertionError error

	for key, value := range j.necessaryFields {
		bodyValue := js.Get(key).Value()

		if !reflect.DeepEqual(bodyValue, value) {
			assertionError = fmt.Errorf(`json assertion failed for "%s" field: expected "%v", got "%v"`, key, value, bodyValue)
			break
		}
	}

	return assertionError
}

func (j *subJSON) Log(t testing.TB) {
	t.Log("Testing request for a required json fields")
}

func (j *subJSON) Error(t testing.TB, err error) {
	t.Errorf(assertErrorTemplate, err)
}
