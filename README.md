# go-api-router
> Lightweight API [httprouter](https://github.com/julienschmidt/httprouter) middleware: cors, logging, and standardized error handling. Extends Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) package.

[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/go-api-router)](https://golang.org/)
[![Build Status](https://travis-ci.com/mrz1836/go-api-router.svg?branch=master)](https://travis-ci.com/mrz1836/go-api-router)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-api-router?style=flat)](https://goreportcard.com/report/github.com/mrz1836/go-api-router)
[![codecov](https://codecov.io/gh/mrz1836/go-api-router/branch/master/graph/badge.svg)](https://codecov.io/gh/mrz1836/go-api-router)
[![Release](https://img.shields.io/github/release-pre/mrz1836/go-api-router.svg?style=flat)](https://github.com/mrz1836/go-api-router/releases)
[![GoDoc](https://godoc.org/github.com/mrz1836/go-api-router?status.svg&style=flat)](https://pkg.go.dev/github.com/mrz1836/go-api-router)

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
```shell script
go get -u github.com/mrz1836/go-api-router
```

## Documentation
You can view the generated [documentation here](https://pkg.go.dev/github.com/mrz1836/go-api-router).

### Features
- Uses the fastest router: Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter)
- Uses gofrs [uuid](https://github.com/gofrs/uuid) package to guarantee unique request ids
- Uses MrZ's [go-logger](https://github.com/mrz1836/go-logger) for either local or remote logging via [Log Entries (Rapid7)](https://www.rapid7.com/products/insightops/)
- Uses MrZ's [go-parameters](https://github.com/mrz1836/go-parameters) for parsing any type of incoming parameter with ease
- Added basic middleware support from Rileyr's [middleware](https://github.com/rileyr/middleware)
- Added Additional CORS Functionality
- Standardized Error Responses for API Requests
- Centralized Logging on All Requests (requesting user info & request time)
- Custom Response Writer for Etag and Cache Support
- `GetClientIPAddress()` safely detects IP addresses behind load balancers
- `GetParams()` parses parameters only once
- `FilterMap()` can remove any confidential parameters from logs

<details>
<summary><strong><code>Package Dependencies</code></strong></summary>

- Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) package.
- Rileyr's [middleware](https://github.com/rileyr/middleware) package.
- gofrs [uuid](https://github.com/gofrs/uuid) package.
- MrZ's [go-logger](https://github.com/mrz1836/go-logger) and [go-parameters](https://github.com/mrz1836/go-parameters) package.
</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to Github and can be installed via: `brew install goreleaser`.

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production.
</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>

View all `makefile` commands
```shell script
make help
```

List of all current commands:
```text
bench                          Run all benchmarks in the Go application
clean                          Remove previous builds and any test cache data
clean-mods                     Remove all the Go mod cache
coverage                       Shows the test coverage
godocs                         Sync the latest tag with GoDocs
help                           Show all make commands available
lint                           Run the Go lint application
release                        Full production release (creates release in Github)
release-test                   Full production test release (everything except deploy)
release-snap                   Test the full release (build binaries)
run-examples                   Runs all the examples
tag                            Generate a new tag and push (IE: make tag version=0.0.0)
tag-remove                     Remove a tag if found (IE: make tag-remove version=0.0.0)
tag-update                     Update an existing tag to current commit (IE: make tag-update version=0.0.0)
test                           Runs vet, lint and ALL tests
test-short                     Runs vet, lint and tests (excludes integration tests)
update                         Update all project dependencies
update-releaser                Update the goreleaser application
vet                            Run the Go vet application
```
</details>

## Examples & Tests
All unit tests and [examples](examples/examples.go) run via [Travis CI](https://travis-ci.com/mrz1836/go-api-router) and uses [Go version 1.14.x](https://golang.org/doc/go1.14). View the [deployment configuration file](.travis.yml).

Run all tests (including integration tests)
```shell script
make test
```

Run tests (excluding integration tests)
```shell script
make test-short
```

View and run the examples:
```shell script
make run-examples
```

## Benchmarks
Run the Go benchmarks:
```shell script
make bench
```

## Code Standards
Read more about this Go project's [code standards](CODE_STANDARDS.md).

## Usage
View the [examples](examples/examples.go)

Basic implementation:
```go
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

This project uses gofrs [uuid](https://github.com/gofrs/uuid) package.

This project uses Rileyr's [middleware](https://github.com/rileyr/middleware) package.

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

Support the development of this project 🙏

[![Donate](https://img.shields.io/badge/donate-bitcoin-brightgreen.svg)](https://mrz1818.com/?tab=tips&af=go-api-router)

## License

![License](https://img.shields.io/github/license/mrz1836/go-api-router.svg?style=flat)
