# go-api-router
**go-api-router** is a lightweight API router middleware: cors, logging, and standardized error handling. This package is intended to be used with Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) and uses MrZ's [go-logger](https://github.com/mrz1836/go-logger).

| | | | | | | |
|-|-|-|-|-|-|-|
| ![License](https://img.shields.io/github/license/mrz1836/go-api-router.svg?style=flat&p=1) | [![Report](https://goreportcard.com/badge/github.com/mrz1836/go-api-router?style=flat&p=1)](https://goreportcard.com/report/github.com/mrz1836/go-api-router)  | [![Codacy Badge](https://api.codacy.com/project/badge/Grade/01708ca3079e4933bafb3b39fe2aaa9d)](https://www.codacy.com/app/mrz1818/go-api-router?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mrz1836/go-api-router&amp;utm_campaign=Badge_Grade) |  [![Build Status](https://travis-ci.com/mrz1836/go-api-router.svg?branch=master)](https://travis-ci.com/mrz1836/go-api-router)   |  [![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat)](https://github.com/RichardLitt/standard-readme) | [![Release](https://img.shields.io/github/release-pre/mrz1836/go-api-router.svg?style=flat)](https://github.com/mrz1836/go-api-router/releases) | [![GoDoc](https://godoc.org/github.com/mrz1836/go-api-router?status.svg&style=flat)](https://godoc.org/github.com/mrz1836/go-api-router) |

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

**go-api-router** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy) and [dep](https://github.com/golang/dep).
```bash
$ go get -u github.com/mrz1836/go-api-router
```

Updating dependencies in **go-api-router**:
```bash
$ cd ../go-api-router
$ dep ensure -update -v
```

### Package Dependencies
- Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) package.
- Satori's [go.uuid](https://github.com/satori/go.uuid) package.
- MrZ's [go-logger](https://github.com/mrz1836/go-logger) package.

## Documentation
You can view the generated [documentation here](https://godoc.org/github.com/mrz1836/go-api-router).

### Features


## Examples & Tests
All unit tests and [examples](middleware_test.go) run via [Travis CI](https://travis-ci.com/mrz1836/go-api-router) and uses [Go version 1.12.x](https://golang.org/doc/go1.12). View the [deployment configuration file](.travis.yml).

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

## Benchmarks
Run the Go [benchmarks](pipl_test.go):
```bash
$ cd ../go-api-router
$ go test -bench . -benchmem
```

## Code Standards
Read more about this Go project's [code standards](CODE_STANDARDS.md).

## Usage

Basic implementation:
```golang
package main

import (

)

func main() {

}
```

## Maintainers

[@MrZ1836](https://github.com/mrz1836)

## Contributing

This project uses Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter) package.

This project uses Satori's [go.uuid](https://github.com/satori/go.uuid) package.

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

Support the development of this project 🙏

[![Donate](https://img.shields.io/badge/donate-bitcoin-brightgreen.svg)](https://mrz1818.com/?tab=tips&af=go-api-router)

## License

![License](https://img.shields.io/github/license/mrz1836/go-api-router.svg?style=flat&p=1)
