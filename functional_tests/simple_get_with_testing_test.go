// nolint dupl
package functional_tests

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/maxcnunes/httpfake"
)

// TestSimpleGet tests a fake server handling a GET request
func TestSimpleGetWithTesting(t *testing.T) {
	fakeService := httpfake.New(httpfake.WithTesting(t))
	defer fakeService.Close()

	// register a handler for our fake service
	fakeService.NewHandler().
		Get("/users?movie=dreamers").
		AssertQueryValue("movie", "dreamers").
		Reply(200).
		BodyString(`[{"username": "dreamer","movie": "dreamers"}]`)

	res, err := http.Get(fakeService.ResolveURL("/users?movie=dreamers"))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close() // nolint errcheck

	// Check the status code is what we expect
	if status := res.StatusCode; status != 200 {
		t.Errorf("request returned wrong status code: got %v want %v",
			status, 200)
	}

	// Check the response body is what we expect
	expected := `[{"username": "dreamer","movie": "dreamers"}]`
	body, _ := ioutil.ReadAll(res.Body)
	if bodyString := string(body); bodyString != expected {
		t.Errorf("request returned unexpected body: got %v want %v",
			bodyString, expected)
	}
}
