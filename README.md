# go-api-router
> Lightweight API [httprouter](https://github.com/julienschmidt/httprouter) middleware: cors, logging, and standardized error handling.

[![Release](https://img.shields.io/github/release-pre/mrz1836/go-api-router.svg?logo=github&style=flat&v=3)](https://github.com/mrz1836/go-api-router/releases)
[![Build Status](https://img.shields.io/github/workflow/status/mrz1836/go-api-router/run-go-tests?logo=github&v=3)](https://github.com/mrz1836/go-api-router/actions)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-api-router?style=flat&v=3)](https://goreportcard.com/report/github.com/mrz1836/go-api-router)
[![codecov](https://codecov.io/gh/mrz1836/go-api-router/branch/master/graph/badge.svg?v=3)](https://codecov.io/gh/mrz1836/go-api-router)
[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/go-api-router?v=3)](https://golang.org/)
<br>
[![Mergify Status](https://img.shields.io/endpoint.svg?url=https://api.mergify.com/v1/badges/mrz1836/go-api-router&style=flat&v=1)](https://mergify.io)
[![Sponsor](https://img.shields.io/badge/sponsor-MrZ-181717.svg?logo=github&style=flat&v=3)](https://github.com/sponsors/mrz1836)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat&v=3)](https://mrz1818.com/?tab=tips&utm_source=github&utm_medium=sponsor-link&utm_campaign=go-api-router&utm_term=go-api-router&utm_content=go-api-router)

<br/>

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

<br/>

## Installation

**go-api-router** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/mrz1836/go-api-router
```

<br/>

## Documentation
View the generated [documentation](https://pkg.go.dev/github.com/mrz1836/go-api-router)

[![GoDoc](https://godoc.org/github.com/mrz1836/go-api-router?status.svg&style=flat)](https://pkg.go.dev/github.com/mrz1836/go-api-router)

### Features
- Uses the fastest router: Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter)
- Uses gofr's [uuid](https://github.com/gofrs/uuid) package to guarantee unique request ids
- Uses MrZ's [go-logger](https://github.com/mrz1836/go-logger) for either local or remote logging via [Log Entries (Rapid7)](https://www.rapid7.com/products/insightops/)
- Uses MrZ's [go-parameters](https://github.com/mrz1836/go-parameters) for parsing any type of incoming parameter with ease
- Optional: [NewRelic](https://docs.newrelic.com/docs/agents/go-agent/get-started/go-agent-compatibility-requirements/) support!
- Added basic middleware support from Rileyr's [middleware](https://github.com/rileyr/middleware)
- Added additional CORS functionality
- Standardized error responses for API requests
- Centralized logging on all requests (requesting user info & request time)
- Custom response writer for Etag and cache support
- `GetClientIPAddress()` safely detects IP addresses behind load balancers
- `GetParams()` parses parameters only once
- `FilterMap()` removes any confidential parameters from logs
- ...and more!

<details>
<summary><strong><code>Package Dependencies</code></strong></summary>
<br/>

- Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) package
- Rileyr's [middleware](https://github.com/rileyr/middleware) package
- gofrs [uuid](https://github.com/gofrs/uuid) package
- MrZ's [go-logger](https://github.com/mrz1836/go-logger) and [go-parameters](https://github.com/mrz1836/go-parameters) package
- NewRelic's [go-agent](https://github.com/newrelic/go-agent)
</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to Github and can be installed via: `brew install goreleaser`.

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production.
</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>
<br/>

View all `makefile` commands
```shell script
make help
```

List of all current commands:
```text
all                  Runs lint, test and vet
clean                Remove previous builds and any test cache data
clean-mods           Remove all the Go mod cache
coverage             Shows the test coverage
godocs               Sync the latest tag with GoDocs
help                 Show this help message
install              Install the application
install-go           Install the application (Using Native Go)
lint                 Run the golangci-lint application (install if not found)
release              Full production release (creates release in Github)
release              Runs common.release then runs godocs
release-snap         Test the full release (build binaries)
release-test         Full production test release (everything except deploy)
replace-version      Replaces the version in HTML/JS (pre-deploy)
run-examples         Runs all the examples
tag                  Generate a new tag and push (tag version=0.0.0)
tag-remove           Remove a tag if found (tag-remove version=0.0.0)
tag-update           Update an existing tag to current commit (tag-update version=0.0.0)
test                 Runs vet, lint and ALL tests
test-ci              Runs all tests via CI (exports coverage)
test-ci-no-race      Runs all tests via CI (no race) (exports coverage)
test-ci-short        Runs unit tests via CI (exports coverage)
test-short           Runs vet, lint and tests (excludes integration tests)
uninstall            Uninstall the application (and remove files)
update-linter        Update the golangci-lint package (macOS only)
vet                  Run the Go vet application
```
</details>
 
<br/>

## Examples & Tests
All unit tests and [examples](examples) run via [Github Actions](https://github.com/mrz1836/go-api-router/actions) and
uses [Go version 1.16.x](https://golang.org/doc/go1.16). View the [configuration file](.github/workflows/run-tests.yml).

Run all tests (including integration tests)
```shell script
make test
```

Run tests (excluding integration tests)
```shell script
make test-short
```

Run the examples:
```shell script
make run-examples
```

<br/>

## Benchmarks
Run the Go benchmarks:
```shell script
make bench
```

<br/>

## Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

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
	_, _ = fmt.Fprint(w, "This is a simple route example!")
}
```

<br/>

## Maintainers
| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:------------------------------------------------------------------------------------------------:|
|                                [MrZ](https://github.com/mrz1836)                                 |


<br/>

## Contributing
View the [contributing guidelines](.github/CONTRIBUTING.md) and please follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:! 
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:. 
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/mrz1836) :clap: 
or by making a [**bitcoin donation**](https://mrz1818.com/?tab=tips&utm_source=github&utm_medium=sponsor-link&utm_campaign=go-api-router&utm_term=go-api-router&utm_content=go-api-router) to ensure this journey continues indefinitely! :rocket:

[![Stars](https://img.shields.io/github/stars/mrz1836/go-api-router?label=Please%20like%20us&style=social)](https://github.com/mrz1836/go-api-router/stargazers)

<br/>

## License
[![License](https://img.shields.io/github/license/mrz1836/go-api-router.svg?style=flat)](LICENSE)
