# Contributing to httpfake

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

There are some ways of contributing to httpfake

* Report an issue.
* Contribute to the code base.

## Report an issue

* Before opening the issue make sure there isn't an issue opened for the same problem
* Include the Go and httpfake version you are using
* If is a bug, please include an example to reproduce the problem ([example](/functional_tests))

## Contribute to the code base

### Pull Request

* Please discuss the suggested changes on a issue before working on it. Just to make sure the change makes sense before you spending any time on it.
* Include a [functional test](/functional_tests) for new features.

### Running lint

```
make lint
```

### Running tests

All tests:
```
make test
```

A single test:
```
go test -run TestSimplePost ./...
```

