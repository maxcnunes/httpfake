// nolint dupl
package functional_tests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/maxcnunes/httpfake"
)

// TestSimplePostWithTesting tests a fake server handling a POST request
func TestSimplePostWithTesting(t *testing.T) {
	fakeService := httpfake.New(httpfake.WithTesting(t))
	defer fakeService.Server.Close()

	// register a handler for our fake service
	fakeService.NewHandler().
		Post("/users").
		AssertHeaders("Authorization").
		AssertHeaderValue("Content-Type", "application/json").
		Reply(201).
		BodyString(`{"id": 1, "username": "dreamer"}`)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	sendBody := bytes.NewBuffer([]byte(`{"username": "dreamer"}`))
	req, err := http.NewRequest(http.MethodPost, fakeService.ResolveURL("/users"), sendBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer some-token")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
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
