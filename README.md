httpfake
========

[![LICENSE](https://img.shields.io/badge/license-MIT-orange.svg)](LICENSE)
[![Godocs](https://img.shields.io/badge/golang-documentation-blue.svg)](https://godoc.org/github.com/maxcnunes/httpfake)
[![Build Status](https://travis-ci.org/maxcnunes/httpfake.svg?branch=master)](https://travis-ci.org/maxcnunes/httpfake)
[![Coverage Status](https://coveralls.io/repos/github/maxcnunes/httpfake/badge.svg?branch=master)](https://coveralls.io/github/maxcnunes/httpfake?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxcnunes/httpfake)](https://goreportcard.com/report/github.com/maxcnunes/httpfake)

httpfake provides is a simple wrapper for [httptest](https://golang.org/pkg/net/http/httptest/) with a handful chainable API for setting up handlers to a fake server. This package is aimed to be used in tests where the original external server must not be reached. Instead is used in its place a fake server which can be configured to handle any request as desired.

## Installation

```
go get -u github.com/maxcnunes/httpfake
```

or

```
govendor fetch github.com/maxcnunes/httpfake
```

> If possible give preference for using vendor. This way the version is locked up as a dependency in your project.

## Changelog

See [Releases](https://github.com/maxcnunes/httpfake/releases) for detailed history changes.

## API

See [godoc reference](https://godoc.org/github.com/maxcnunes/httpfake) for detailed API documentation.

## Examples

For a full list of examples please check out the [functional_tests folder](/functional_tests).

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
res, err := http.Get(fakeService.ResolveURL("/users"))
```

## Contributing

See the [Contributing guide](/CONTRIBUTING.md) for steps on how to contribute to this project.

## Reference

This package was heavily inspired by [gock](https://github.com/h2non/gock). Check that you if you prefer mocking your requests.
