.PHONY: deps deps-integration lint test test-integration ci

deps:
	go get -v -t -d ./...

deps-integration:
	go get -v -t -d -tags integration ./...

lint:
	golangci-lint run -D=lll,gochecknoglobals,gosec,goconst,gocritic

test:
	go test -v -race ./...

test-integration:
	go test -v -race -tags integration -onfidoToken=$(ONFIDO_TOKEN)

ci: deps lint test
