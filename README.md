# go-api-router
**go-api-router** is a lightweight API [httprouter](https://github.com/julienschmidt/httprouter) middleware: cors, logging, and standardized error handling. Extends Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) package.

[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/go-api-router)](https://golang.org/)
[![Build Status](https://travis-ci.com/mrz1836/go-api-router.svg?branch=master)](https://travis-ci.com/mrz1836/go-api-router)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-api-router?style=flat)](https://goreportcard.com/report/github.com/mrz1836/go-api-router)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/0b377a0d1dde4b6ba189545aa7ee2e17)](https://www.codacy.com/app/mrz1818/go-api-router?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mrz1836/go-api-router&amp;utm_campaign=Badge_Grade)
[![Release](https://img.shields.io/github/release-pre/mrz1836/go-api-router.svg?style=flat)](https://github.com/mrz1836/go-api-router/releases)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat)](https://github.com/RichardLitt/standard-readme)
[![GoDoc](https://godoc.org/github.com/mrz1836/go-api-router?status.svg&style=flat)](https://godoc.org/github.com/mrz1836/go-api-router)

## Table of Contents
- [Installation](#installation)
- [Documentation](#documentation)
- [Examples & Tests](#examples--tests)
- [Benchmarks](#benchmarks)
- [Code Standards](#code-standards)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Installation

**go-api-router** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```bash
$ go get -u github.com/mrz1836/go-api-router
```

### Package Dependencies
- Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) package.
- Rileyr's [middleware](https://github.com/rileyr/middleware) package.
- Satori's [go.uuid](https://github.com/satori/go.uuid) package.
- MrZ's [go-logger](https://github.com/mrz1836/go-logger) and [go-parameters](https://github.com/mrz1836/go-parameters) package.

## Documentation
You can view the generated [documentation here](https://godoc.org/github.com/mrz1836/go-api-router).

### Features
- Uses the fastest router: Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter)
- Uses Satori's [go.uuid](https://github.com/satori/go.uuid) package to guarantee unique request ids
- Uses MrZ's [go-logger](https://github.com/mrz1836/go-logger) for either local or remote logging via [LogEntries](https://logentries.com/)
- Uses MrZ's [go-parameters](https://github.com/mrz1836/go-parameters) for parsing any type of incoming parameter with ease
- Added basic middleware support from Rileyr's [middleware](https://github.com/rileyr/middleware)
- Added Additional CORS Functionality
- Standardized Error Responses for API Requests
- Centralized Logging on All Requests (requesting user info & request time)
- Custom Response Writer for Etag and Cache Support
- `GetClientIPAddress()` safely detects IP addresses behind load balancers
- `GetParams()` parses parameters only once
- `FilterMap()` can remove any confidential parameters from logs

## Examples & Tests
All unit tests and [examples](examples/examples.go) run via [Travis CI](https://travis-ci.com/mrz1836/go-api-router) and uses [Go version 1.13.x](https://golang.org/doc/go1.13). View the [deployment configuration file](.travis.yml).

Run all tests (including integration tests)
```bash
$ cd ../go-api-router
$ go test ./... -v
```

Run tests (excluding integration tests)
```bash
$ cd ../go-api-router
$ go test ./... -v -test.short
```

View and run the examples:
```bash
$ cd ../go-api-router/examples
$ go run examples.go
```

## Benchmarks
Run the Go benchmarks:
```bash
$ cd ../go-api-router
$ go test -bench . -benchmem
```

## Code Standards
Read more about this Go project's [code standards](CODE_STANDARDS.md).

## Usage
View the [examples](examples/examples.go)

Basic implementation:
```golang
package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mrz1836/go-api-router"
	"github.com/mrz1836/go-logger"
)

func main() {
	// Load the router & middleware
	router := apirouter.New()

	// Set the main index page (navigating to slash)
	router.HTTPRouter.GET("/", router.Request(index))

	// Serve the router!
	logger.Fatalln(http.ListenAndServe(":3000", router.HTTPRouter))
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprint(w, "This is a simple API example!")
}
```

## Maintainers

| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:---:|
| [MrZ](https://github.com/mrz1836) |


## Contributing

This project uses Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) package.

This project uses Satori's [go.uuid](https://github.com/satori/go.uuid) package.

This project uses Rileyr's [middleware](https://github.com/rileyr/middleware) package.

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

Support the development of this project 🙏

[![Donate](https://img.shields.io/badge/donate-bitcoin-brightgreen.svg)](https://mrz1818.com/?tab=tips&af=go-api-router)

## License

![License](https://img.shields.io/github/license/mrz1836/go-api-router.svg?style=flat)
