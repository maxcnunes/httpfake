// nolint dupl
package functional_tests

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/voronelf/httpfake"
)

// TestNotRegisteredMethod tests a fake server handling a GET request
// for a not registered method
func TestNotRegisteredMethod(t *testing.T) {
	fakeService := httpfake.New()
	defer fakeService.Server.Close()

	// register a handler for our fake service
	fakeService.NewHandler().
		Post("/users").
		Reply(200).
		BodyString(`[{"username": "dreamer"}]`)

	res, err := http.Get(fakeService.ResolveURL("/clients"))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close() // nolint errcheck

	// Check the status code is what we expect
	if status := res.StatusCode; status != 404 {
		t.Errorf("request returned wrong status code: got %v want %v",
			status, 404)
	}

	// Check the response body is what we expect
	expected := ""
	body, _ := ioutil.ReadAll(res.Body)
	if bodyString := string(body); bodyString != expected {
		t.Errorf("request returned unexpected body: got %v want %v",
			bodyString, expected)
	}
}
