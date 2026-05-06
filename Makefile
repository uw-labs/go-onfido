GOLANGCI_LINT_VERSION ?= v1.16.0

.PHONY: deps deps-integration install-lint lint test test-integration ci

deps:
	go get -v -t -d ./...

deps-integration:
	go get -v -t -d -tags integration ./...

install-lint:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

lint:
	golangci-lint run --enable-all -D=lll,gochecknoglobals,gosec,goconst,gocritic

test:
	go test -v -race ./...

test-integration:
	go test -v -race -tags integration -onfidoToken=$(ONFIDO_TOKEN)

ci: deps lint test
