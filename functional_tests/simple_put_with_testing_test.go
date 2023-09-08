// nolint dupl
package functional_tests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/voronelf/httpfake"
)

// TestSimplePutWithTesting tests a fake server handling a PUT request
func TestSimplePutWithTesting(t *testing.T) {
	fakeService := httpfake.New(httpfake.WithTesting(t))
	defer fakeService.Close()

	// register a handler for our fake service
	fakeService.NewHandler().
		Put("/users/1").
		AssertBody([]byte(`{"username": "dreamer"}`)).
		Reply(200).
		BodyString(`{"id": 1,"username": "dreamer"}`)

	sendBody := bytes.NewBuffer([]byte(`{"username": "dreamer"}`))
	req, err := http.NewRequest("PUT", fakeService.ResolveURL("/users/1"), sendBody)
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
	expected := `{"id": 1,"username": "dreamer"}`
	body, _ := ioutil.ReadAll(res.Body)
	if bodyString := string(body); bodyString != expected {
		t.Errorf("request returned unexpected body: got %v want %v",
			bodyString, expected)
	}
}
