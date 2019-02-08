# go-onfido [![CircleCI](https://circleci.com/gh/uw-labs/go-onfido.svg?style=svg)](https://circleci.com/gh/uw-labs/go-onfido) [![Go Report Card](https://goreportcard.com/badge/github.com/tumelohq/go-onfido)](https://goreportcard.com/report/github.com/tumelohq/go-onfido)

Client for the [Onfido API](https://documentation.onfido.com/)

[![go-doc](https://godoc.org/github.com/tumelohq/go-onfido?status.svg)](https://godoc.org/github.com/tumelohq/go-onfido)

> This library was built for Utility Warehouse internal projects, so priority was given to supporting the
features we needed. If the library is missing a feature from the API, raise an issue or ideally open a PR.

## Installation

To install go-onfido, use `go get`:

```
go get github.com/uw-labs/go-onfio
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

Now checkout some of the [examples](https://github.com/tumelohq/go-onfido/tree/master/examples)


