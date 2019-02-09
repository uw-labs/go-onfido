.PHONY: unit_test
unit_tests:
	- @go test -v -race .

.PHONY: integration_tests
integration_tests:
	- @go test -v -race -tags integration -onfidoToken=${ONFIDO_TOKEN}

.PHONY: fmt
fmt:
	- @go fmt .
