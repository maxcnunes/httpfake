// nolint dupl
package functional_tests

import (
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	"github.com/maxcnunes/httpfake"
)

// TestSimpleGetWithRutines tests a fake server handling a GET request
func TestSimpleGetWithRutines(t *testing.T) {
	numOfRequests := 3
	fakeService := httpfake.New()
	defer fakeService.Server.Close()

	// register a handler for our fake service
	fakeService.NewHandler().
		Get("/users").
		Reply(200).
		BodyString(`[{"username": "dreamer"}]`)

	var wg sync.WaitGroup
	wg.Add(numOfRequests)
	for i := 0; i < numOfRequests; i++ {
		go func(wg *sync.WaitGroup) {
			res, err := http.Get(fakeService.ResolveURL("/users"))
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

			wg.Done()
		}(&wg)
	}

	wg.Wait()
}
