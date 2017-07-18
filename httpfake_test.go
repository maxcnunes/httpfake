package httpfake

import (
	"testing"
)

func TestResolveURL(t *testing.T) {
	fakeService := New()
	defer fakeService.Server.Close()

	testCases := []struct {
		resolved string
		expected string
	}{
		{
			resolved: fakeService.ResolveURL("/users"),
			expected: fakeService.Server.URL + "/users",
		},
		{
			resolved: fakeService.ResolveURL("/users/%v", 1),
			expected: fakeService.Server.URL + "/users/1",
		},
	}

	for _, tc := range testCases {
		t.Run("ResolveURL", func(t *testing.T) {
			if tc.resolved != tc.expected {
				t.Errorf("returned unexpected URL: got %v want %v",
					tc.resolved, tc.expected)
			}
		})
	}

}
