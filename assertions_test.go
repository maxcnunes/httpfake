// nolint dupl gocyclo
package httpfake

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

type mockTester struct {
	testing.TB
	buf *bytes.Buffer
}

func (t *mockTester) Log(args ...interface{}) {
	t.buf.WriteString(fmt.Sprintln(args...))
}

func (t *mockTester) Logf(format string, args ...interface{}) {
	t.buf.WriteString(fmt.Sprintf(format, args...))
}

func (t *mockTester) Errorf(format string, args ...interface{}) {
	t.buf.WriteString(fmt.Sprintf(format, args...))
}

func TestAssertors_Assert(t *testing.T) {
	tests := []struct {
		name           string
		assertor       Assertor
		requestBuilder func() (*http.Request, error)
		expectedErr    string
	}{
		{
			name: "RequiredHeaders should return no error with a proper request",
			assertor: &RequiredHeaders{
				Keys: []string{"test-header-1", "test-header-2"},
			},
			requestBuilder: func() (*http.Request, error) {
				testReq, err := http.NewRequest(http.MethodPost, "http://fake.url", nil)
				if err != nil {
					return nil, err
				}

				testReq.Header.Set("test-header-1", "mock-value-1")
				testReq.Header.Set("test-header-2", "mock-value-2")

				return testReq, nil
			},
			expectedErr: "",
		},
		{
			name: "RequiredHeaders should return an error if a request is missing a required header",
			assertor: &RequiredHeaders{
				Keys: []string{"test-header-1", "test-header-2"},
			},
			requestBuilder: func() (*http.Request, error) {
				testReq, err := http.NewRequest(http.MethodPost, "http://fake.url", nil)
				if err != nil {
					return nil, err
				}

				testReq.Header.Set("test-header-2", "mock-value-2")

				return testReq, nil
			},
			expectedErr: "missing required header(s): test-header-1",
		},
		{
			name: "RequiredHeaderValue should return no error with a proper request",
			assertor: &RequiredHeaderValue{
				Key:           "test-header-1",
				ExpectedValue: "mock-value-1",
			},
			requestBuilder: func() (*http.Request, error) {
				testReq, err := http.NewRequest(http.MethodPost, "http://fake.url", nil)
				if err != nil {
					return nil, err
				}

				testReq.Header.Set("test-header-1", "mock-value-1")

				return testReq, nil
			},
			expectedErr: "",
		},
		{
			name: "RequiredHeaderValue should return an error if a request is missing a required header",
			assertor: &RequiredHeaderValue{
				Key:           "test-header-1",
				ExpectedValue: "mock-value-1",
			},
			requestBuilder: func() (*http.Request, error) {
				testReq, err := http.NewRequest(http.MethodPost, "http://fake.url", nil)
				if err != nil {
					return nil, err
				}

				return testReq, nil
			},
			expectedErr: "header test-header-1 does not have the expected value; expected  to equal mock-value-1",
		},
		{
			name: "RequiredQueries should return no error with a proper request",
			assertor: &RequiredQueries{
				Keys: []string{"query-1", "query-2"},
			},
			requestBuilder: func() (*http.Request, error) {
				u := "http://fake.url?query-1=apples&query-2=oranges"
				mockURL, err := url.Parse(u)
				if err != nil {
					return nil, err
				}

				testReq, err := http.NewRequest(http.MethodPost, mockURL.Host, nil)
				if err != nil {
					return nil, err
				}
				testReq.URL = mockURL

				return testReq, nil
			},
			expectedErr: "",
		},
		{
			name: "RequiredQueries should return an error if a request is missing the a required query params",
			assertor: &RequiredQueries{
				Keys: []string{"query-1", "query-3"},
			},
			requestBuilder: func() (*http.Request, error) {
				u := "http://fake.url?query-2=oranges"
				mockURL, err := url.Parse(u)
				if err != nil {
					return nil, err
				}

				testReq, err := http.NewRequest(http.MethodPost, mockURL.Host, nil)
				if err != nil {
					return nil, err
				}
				testReq.URL = mockURL

				return testReq, nil
			},
			expectedErr: "missing required query parameter(s): query-1, query-3",
		},
		{
			name: "RequiredQueryValue should return no error with a proper request",
			assertor: &RequiredQueryValue{
				Key:           "query-1",
				ExpectedValue: "value-1",
			},
			requestBuilder: func() (*http.Request, error) {
				u := "http://fake.url?query-1=value-1"
				mockURL, err := url.Parse(u)
				if err != nil {
					return nil, err
				}

				testReq, err := http.NewRequest(http.MethodPost, mockURL.Host, nil)
				if err != nil {
					return nil, err
				}
				testReq.URL = mockURL

				return testReq, nil
			},
			expectedErr: "",
		},
		{
			name: "RequiredQueryValue should return an error if a request is missing the a required query param",
			assertor: &RequiredQueryValue{
				Key:           "query-1",
				ExpectedValue: "apples",
			},
			requestBuilder: func() (*http.Request, error) {
				u := "http://fake.url"
				mockURL, err := url.Parse(u)
				if err != nil {
					return nil, err
				}

				testReq, err := http.NewRequest(http.MethodPost, mockURL.Host, nil)
				if err != nil {
					return nil, err
				}
				testReq.URL = mockURL

				return testReq, nil
			},
			expectedErr: "query query-1 does not have the expected value; expected  to equal apples",
		},
		{
			name: "RequiredQueryValue should return an error if a request has an incorrect query param value",
			assertor: &RequiredQueryValue{
				Key:           "query-1",
				ExpectedValue: "apples",
			},
			requestBuilder: func() (*http.Request, error) {
				u := "http://fake.url?query-1=oranges"
				mockURL, err := url.Parse(u)
				if err != nil {
					return nil, err
				}

				testReq, err := http.NewRequest(http.MethodPost, mockURL.Host, nil)
				if err != nil {
					return nil, err
				}
				testReq.URL = mockURL

				return testReq, nil
			},
			expectedErr: "query query-1 does not have the expected value; expected oranges to equal apples",
		},
		{
			name: "RequiredBody should return no error with a proper request",
			assertor: &RequiredBody{
				ExpectedBody: []byte(`{"testObj": {"data": {"fakeData": "testdata"}}}`),
			},
			requestBuilder: func() (*http.Request, error) {
				reader := bytes.NewBuffer([]byte(`{"testObj": {"data": {"fakeData": "testdata"}}}`))

				testReq, err := http.NewRequest(http.MethodPost, "http://fake.url", reader)
				if err != nil {
					return nil, err
				}

				return testReq, nil
			},
			expectedErr: "",
		},
		{
			name: "RequiredBody should return an error if the body is not what's expected",
			assertor: &RequiredBody{
				ExpectedBody: []byte(`{"testObj": {"data": {"fakeData": "testdata"}}}`),
			},
			requestBuilder: func() (*http.Request, error) {
				reader := bytes.NewBuffer([]byte(`{"testObj": {"data": {"badData": "bad"}}}`))

				testReq, err := http.NewRequest(http.MethodPost, "http://fake.url", reader)
				if err != nil {
					return nil, err
				}

				return testReq, nil
			},
			expectedErr: "request body does not have the expected value; expected {\"testObj\": {\"data\": {\"badData\": \"bad\"}}} to equal {\"testObj\": {\"data\": {\"fakeData\": \"testdata\"}}}",
		},
		{
			name: "RequiredBody should handle a nil body without panic",
			assertor: &RequiredBody{
				ExpectedBody: []byte(`{"testObj": {"data": {"fakeData": "testdata"}}}`),
			},
			requestBuilder: func() (*http.Request, error) {

				testReq, err := http.NewRequest(http.MethodPost, "http://fake.url", nil)
				if err != nil {
					return nil, err
				}

				return testReq, nil
			},
			expectedErr: "error reading the request body; the request body is nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testReq, err := tt.requestBuilder()
			if err != nil {
				t.Fatalf("error setting up fake request: %#v", err)
			}

			err = tt.assertor.Assert(testReq)
			if len(tt.expectedErr) > 0 {
				if err == nil {
					t.Errorf("Expected error %s but err was nil", tt.expectedErr)
					return
				}

				if err.Error() != tt.expectedErr {
					t.Errorf("Assert() error = %v, expected error %s", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error = %v", err)
			}
		})
	}
}

func TestAssertors_Log(t *testing.T) {
	tests := []struct {
		name       string
		mockTester *mockTester
		assertor   Assertor
		expected   string
	}{
		{
			name: "RequiredHeaders Log should log the expected output when called",
			mockTester: &mockTester{
				buf: &bytes.Buffer{},
			},
			assertor: &RequiredHeaders{},
			expected: "Testing request for required headers\n",
		},
		{
			name: "RequiredHeaderValue Log should log the expected output when called",
			mockTester: &mockTester{
				buf: &bytes.Buffer{},
			},
			assertor: &RequiredHeaderValue{Key: "test-key", ExpectedValue: "test-value"},
			expected: "Testing request for required header value [test-key: test-value]",
		},
		{
			name: "RequiredQueries Log should log the expected output when called",
			mockTester: &mockTester{
				buf: &bytes.Buffer{},
			},
			assertor: &RequiredQueries{},
			expected: "Testing request for required query parameters\n",
		},
		{
			name: "RequiredQueryValue Log should log the expected output when called",
			mockTester: &mockTester{
				buf: &bytes.Buffer{},
			},
			assertor: &RequiredQueryValue{Key: "test-key", ExpectedValue: "test-value"},
			expected: "Testing request for required query parameter value [test-key: test-value]",
		},
		{
			name: "RequiredBody Log should log the expected output when called",
			mockTester: &mockTester{
				buf: &bytes.Buffer{},
			},
			assertor: &RequiredBody{},
			expected: "Testing request for required a required body\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertor.Log(tt.mockTester)

			result := tt.mockTester.buf.String()
			if result != tt.expected {
				t.Errorf("Expected Log %s, actual %#v", tt.expected, result)
			}
		})
	}
}

func TestAssertors_Error(t *testing.T) {
	tests := []struct {
		name       string
		mockTester *mockTester
		assertor   Assertor
		expected   string
	}{
		{
			name: "RequiredHeaders Log should log the expected output when called",
			mockTester: &mockTester{
				buf: &bytes.Buffer{},
			},
			assertor: &RequiredHeaders{},
			expected: "assertion error: test error",
		},
		{
			name: "RequiredHeaderValue Log should log the expected output when called",
			mockTester: &mockTester{
				buf: &bytes.Buffer{},
			},
			assertor: &RequiredHeaderValue{Key: "test-key", ExpectedValue: "test-value"},
			expected: "assertion error: test error",
		},
		{
			name: "RequiredQueries Log should log the expected output when called",
			mockTester: &mockTester{
				buf: &bytes.Buffer{},
			},
			assertor: &RequiredQueries{},
			expected: "assertion error: test error",
		},
		{
			name: "RequiredQueryValue Log should log the expected output when called",
			mockTester: &mockTester{
				buf: &bytes.Buffer{},
			},
			assertor: &RequiredQueryValue{Key: "test-key", ExpectedValue: "test-value"},
			expected: "assertion error: test error",
		},
		{
			name: "RequiredBody Log should log the expected output when called",
			mockTester: &mockTester{
				buf: &bytes.Buffer{},
			},
			assertor: &RequiredBody{},
			expected: "assertion error: test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testErr := errors.New("test error")
			tt.assertor.Error(tt.mockTester, testErr)

			result := tt.mockTester.buf.String()
			if result != tt.expected {
				t.Errorf("Expected Error %s, actual %#v", tt.expected, result)
			}
		})
	}
}
