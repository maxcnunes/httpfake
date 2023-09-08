// nolint dupl
package functional_tests

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/voronelf/httpfake"
)

// TestHandleCustomResponder tests a fake server handling a GET request
// with a custom responder. It allows full control over the handler.
func TestHandleCustomResponder(t *testing.T) {
	fakeService := httpfake.New()
	defer fakeService.Server.Close()

	// register a handler for our fake service
	fakeService.NewHandler().
		Get("/users").
		Handle(func(w http.ResponseWriter, r *http.Request, rh *httpfake.Request) {
			w.Header().Add("Header-From-Custom-Responder-X", "indeed")
			w.WriteHeader(200)
			w.Write([]byte("Body-From-Custom-Responder-X")) // nolint
		})

	req, err := http.NewRequest("GET", fakeService.ResolveURL("/users"), nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
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
	expected := "Body-From-Custom-Responder-X"
	body, _ := ioutil.ReadAll(res.Body)
	if bodyString := string(body); bodyString != expected {
		t.Errorf("request returned unexpected body: got %v want %v",
			bodyString, expected)
	}

	// Check the response header is what we expect
	expected = "indeed"
	if header := res.Header.Get("Header-From-Custom-Responder-X"); header != expected {
		t.Errorf("request returned unexpected value for header Header-From-Custom-Responder-X: got %v want %v",
			header, expected)
	}
}
