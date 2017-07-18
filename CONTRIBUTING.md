# Contributing to httpfake

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

There are some ways of contributing to httpfake

* Report an issue.
* Contribute to the code base.

## Report an issue

* Before opening the issue make sure there isn't an issue opened for the same problem
* Include the Go and httpfake version you are using
* If is a bug, please include an example to reproduce the problem

## Contribute to the code base

> If you are interested in sending a pull request. Please discuss it on a issue before working on it.
> Just to make sure the change makes sense before you spending any time on it.

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

