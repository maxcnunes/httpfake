// nolint dupl
package functional_tests

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/voronelf/httpfake"
)

// TestGetQueryWithSpecialChars tests a fake server handling a request with special chars in path and query
func TestGetQueryWithSpecialChars(t *testing.T) {
	fakeService := httpfake.New()
	defer fakeService.Server.Close()

	// register a handler for our fake service
	fakeService.NewHandler().
		Get("/users?name=" + url.QueryEscape("Tim Burton")).
		Reply(200).
		BodyString(`[{"username": "dreamer"}]`)

	// register second handler for our fake service
	fakeService.NewHandler().
		Get("/users?name=other").
		Reply(201).
		BodyString(`[{"username": "other"}]`)

	res, err := http.Get(fakeService.ResolveURL("/users?name=%s", url.QueryEscape("Tim Burton")))
	//res, err := http.Get(fakeService.ResolveURL("/users?name=Tim Burton"))
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
	expected := `[{"username": "dreamer"}]`
	body, _ := ioutil.ReadAll(res.Body)
	if bodyString := string(body); bodyString != expected {
		t.Errorf("request returned unexpected body: got %v want %v",
			bodyString, expected)
	}
}
