// nolint dupl
package functional_tests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/voronelf/httpfake"
)

// TestSimplePost tests a fake server handling a POST request
func TestSimplePost(t *testing.T) {
	fakeService := httpfake.New()
	defer fakeService.Server.Close()

	// register a handler for our fake service
	fakeService.NewHandler().
		Post("/users").
		Reply(201).
		BodyString(`{"id": 1, "username": "dreamer"}`)

	sendBody := bytes.NewBuffer([]byte(`{"username": "dreamer"}`))
	res, err := http.Post(fakeService.ResolveURL("/users"), "application/json", sendBody)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close() // nolint errcheck

	// Check the status code is what we expect
	if status := res.StatusCode; status != 201 {
		t.Errorf("request returned wrong status code: got %v want %v",
			status, 201)
	}

	// Check the response body is what we expect
	expected := `{"id": 1, "username": "dreamer"}`
	body, _ := ioutil.ReadAll(res.Body)
	if bodyString := string(body); bodyString != expected {
		t.Errorf("request returned unexpected body: got %v want %v",
			bodyString, expected)
	}
}
