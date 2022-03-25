[![burrow-test](https://github.com/bdshroyer/burrow/actions/workflows/burrow-test.yml/badge.svg?branch=master)](https://github.com/bdshroyer/burrow/actions/workflows/burrow-test.yml)


# burrow

Golang code for exploring some network centrality ideas. Still very much a WIP.

### Compatibility notes:
I'm currently trying out Go generics in some probability distribution functionality. This means that burrow code may not work with Go versions older than 1.18.

### Important dependencies:
* [Gonum](https://github.com/gonum/gonum): burrow is being built according to gonum Graph interface specs to provide access to the full range of Gonum graph algorithms.
* [Gomega](https://github.com/onsi/gomega): Matchers for testing.
* [Ginkgo](https://github.com/onsi/ginkgo): BDD testing framework for Go.
* [Go/x/exp](https://pkg.go.dev/golang.org/x/exp): Provides the `constraints` package used for generics and type constraints.

### Testing

Burrow tests are written using [Ginkgo](https://onsi.github.io/ginkgo), which runs on top of Go's native testing framework. To execute, run `ginkgo -r` or `go test ./...`.
