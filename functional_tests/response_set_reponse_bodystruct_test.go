// nolint dupl
package functional_tests

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/voronelf/httpfake"
)

// TestResponseBodyStruct tests a fake server handling a GET request
// and responding with a provided struct data
func TestResponseBodyStruct(t *testing.T) {
	fakeService := httpfake.New()
	defer fakeService.Server.Close()

	type user struct {
		UserName string `json:"username"`
	}
	// register a handler for our fake service
	fakeService.NewHandler().
		Get("/users").
		Reply(200).
		BodyStruct(&user{UserName: "dreamer"})

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
	expected := `{"username":"dreamer"}`
	body, _ := ioutil.ReadAll(res.Body)
	if bodyString := string(body); bodyString != expected {
		t.Errorf("request returned unexpected body: got %v want %v",
			bodyString, expected)
	}
}
