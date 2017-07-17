package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/maxcnunes/httpfake"
)

// TestSimplePost tests a fake server handling a POST request
func TestSimplePost(t *testing.T) {
	fakeService := httpfake.New()
	defer fakeService.Server.Close()

	// register a handler for our fake service
	fakeService.NewHandler().
		Post("/users").
		Reply(201).
		BodyString(`{"username": "dreamer"}`)

	res, err := http.Get(fmt.Sprintf("%s/users", fakeService.Server.URL))
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
	expected := `{"username": "dreamer"}`
	body, _ := ioutil.ReadAll(res.Body)
	if bodyString := string(body); bodyString != expected {
		t.Errorf("request returned unexpected body: got %v want %v",
			bodyString, expected)
	}
}
