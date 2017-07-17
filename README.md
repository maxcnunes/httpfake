httpfake
========

*IMPORTANT: In Development*

httpfake is a simple wrapper for [httptest](https://golang.org/pkg/net/http/httptest/) with a handful chainable API for defining handlers to the fake server. This package is aimed to be used in tests where an original external server must not be reached. Instead is used in its place a fake server which can be configured to handle any request as desired.

## Installation

```
go get -u github.com/maxcnunes/httpfake
```

or

```
govendor fetch github.com/maxcnunes/httpfake
```

> If possible give preference for using vendor. This way the version is locked up as a dependency in your project.

## Examples

For a full list of examples please checkout the [examples folder](_examples).

```go
// initialize the faker server
// will bring up a httptest.Server
fakeService := httpfake.New()

// bring down the server once we
// finish running our tests
defer fakeService.Server.Close()

// register a handler for our fake service
fakeService.NewHandler().
  Get("/users").
  Reply(200).
  BodyString(`[{"username": "dreamer"}]`)

// run a real http request to that server
res, err := http.Get(fmt.Sprintf("%s/users", fakeService.Server.URL))
```

## Reference

This package was heavily inspired by [gock](https://github.com/h2non/gock). Check that you if you prefer mocking your requests.
