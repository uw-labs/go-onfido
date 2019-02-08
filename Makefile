.PHONY:
unit_test:
	- @go test -v -race .

fmt:
	- @go fmt .