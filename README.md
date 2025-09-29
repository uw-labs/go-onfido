# go-onfido [![CircleCI](https://circleci.com/gh/uw-labs/go-onfido.svg?style=svg)](https://circleci.com/gh/uw-labs/go-onfido) [![Go Report Card](https://goreportcard.com/badge/github.com/uw-labs/go-onfido)](https://goreportcard.com/report/github.com/uw-labs/go-onfido)

Client for the [Onfido API](https://documentation.onfido.com/)

[![go-doc](https://godoc.org/github.com/uw-labs/go-onfido?status.svg)](https://godoc.org/github.com/uw-labs/go-onfido)

> This library was built for Utility Warehouse internal projects, so priority was given to supporting the
features we needed. If the library is missing a feature from the API, raise an issue or ideally open a PR, 
however please understand that this library is not expected to recieve any ongoing support unless required by
Utility Warehouse.

## Installation

To install go-onfido, use `go get`:

```
go get github.com/uw-labs/go-onfido
```

## Usage

First you're going to need to instantiate a client (grab your [sandbox API key](https://onfido.com/dashboard/v2/#/api/tokens))

```golang
client := onfido.NewClient("test_123")
```

Or you can instantiate usign the env variable `ONFIDO_TOKEN`

```golang
client, err := onfido.NewClientFromEnv()
```

Now checkout some of the [examples](https://github.com/uw-labs/go-onfido/tree/master/examples)


