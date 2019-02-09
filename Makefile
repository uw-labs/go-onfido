.PHONY: unit_test
unit_tests:
	- @go test -v -race .

.PHONY: integration_tests
integration_tests:
	- @go test -v -race -tags integration -onfidoToken=${ONFIDO_TOKEN}

.PHONY: fmt
fmt:
	- @go fmt .

.PHONY: go_get
go_get:
	- @go get .

.PHONY: go_update
go_update:
	- @go get -u .