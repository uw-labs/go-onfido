GOLANGCI_LINT_VERSION ?= v2.12.2

.PHONY: deps deps-integration install-lint lint test test-integration ci

build:
	go build ./...

deps:
	go get -v -t ./...

deps-integration:
	go get -v -t -tags integration ./...

install-lint:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

lint:
	golangci-lint run

test:
	go test -v -race ./...

test-integration:
	go test -v -race -tags integration -onfidoToken=$(ONFIDO_TOKEN)

ci: deps lint test
