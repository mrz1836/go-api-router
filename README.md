# go-api-middleware
**go-api-middleware** is a lightweight middleware for logging, error handling and custom response writer.

| | | | | | | |
|-|-|-|-|-|-|-|
| ![License](https://img.shields.io/github/license/mrz1836/go-api-middleware.svg?style=flat&p=1) | [![Report](https://goreportcard.com/badge/github.com/mrz1836/go-api-middleware?style=flat&p=1)](https://goreportcard.com/report/github.com/mrz1836/go-api-middleware)  | [![Codacy Badge](https://api.codacy.com/project/badge/Grade/01708ca3079e4933bafb3b39fe2aaa9d)](https://www.codacy.com/app/mrz1818/go-api-middleware?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mrz1836/go-api-middleware&amp;utm_campaign=Badge_Grade) |  [![Build Status](https://travis-ci.com/mrz1836/go-api-middleware.svg?branch=master)](https://travis-ci.com/mrz1836/go-api-middleware)   |  [![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat)](https://github.com/RichardLitt/standard-readme) | [![Release](https://img.shields.io/github/release-pre/mrz1836/go-api-middleware.svg?style=flat)](https://github.com/mrz1836/go-api-middleware/releases) | [![GoDoc](https://godoc.org/github.com/mrz1836/go-api-middleware?status.svg&style=flat)](https://godoc.org/github.com/mrz1836/go-api-middleware) |

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

**go-api-middleware** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy) and [dep](https://github.com/golang/dep).
```bash
$ go get -u github.com/mrz1836/go-api-middleware
```

Updating dependencies in **go-api-middleware**:
```bash
$ cd ../go-api-middleware
$ dep ensure -update -v
```

## Documentation
You can view the generated [documentation here](https://godoc.org/github.com/mrz1836/go-api-middleware).

### Features


## Examples & Tests
All unit tests and [examples](middleware_test.go) run via [Travis CI](https://travis-ci.com/mrz1836/go-api-middleware) and uses [Go version 1.12.x](https://golang.org/doc/go1.12). View the [deployment configuration file](.travis.yml).

Run all tests (including integration tests)
```bash
$ cd ../go-api-middleware
$ go test ./... -v
```

Run tests (excluding integration tests)
```bash
$ cd ../go-api-middleware
$ go test ./... -v -test.short
```

## Benchmarks
Run the Go [benchmarks](pipl_test.go):
```bash
$ cd ../go-api-middleware
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

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

Support the development of this project 🙏

[![Donate](https://img.shields.io/badge/donate-bitcoin-brightgreen.svg)](https://mrz1818.com/?tab=tips&af=go-api-middleware)

## License

![License](https://img.shields.io/github/license/mrz1836/go-api-middleware.svg?style=flat&p=1)
