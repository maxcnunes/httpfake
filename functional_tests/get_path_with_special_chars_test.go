// nolint dupl
package functional_tests

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/voronelf/httpfake"
)

// TestGetPathWithSpecialChars tests a fake server handling a request with special chars in path and query
func TestGetPathWithSpecialChars(t *testing.T) {
	fakeService := httpfake.New()
	defer fakeService.Server.Close()

	// register a handler for our fake service
	fakeService.NewHandler().
		Get("/user/+79998887766").
		Reply(200).
		BodyString(`[{"username": "dreamer"}]`)

	res, err := http.Get(fakeService.ResolveURL("/user/+79998887766"))
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
