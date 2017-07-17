TEST_PACKAGES := $(shell go list ./... && go list ./_examples)

install_gometalinter:
	@go get -v github.com/alecthomas/gometalinter
	@gometalinter --install

update_gometalinter:
	@go get -v -u github.com/alecthomas/gometalinter
	@gometalinter --install --update

## lint: Validate golang code
lint: install_gometalinter
	@gometalinter \
		--deadline=120s \
		--line-length=120 \
		--enable-all \
		--vendor ./...

## Perform all tests
test: test/unit

## Perform unit tests
test/unit:
	@go test -v -cover $(TEST_PACKAGES)
