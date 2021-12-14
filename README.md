[![burrow-test](https://github.com/bdshroyer/burrow/actions/workflows/burrow-test.yml/badge.svg?branch=master)](https://github.com/bdshroyer/burrow/actions/workflows/burrow-test.yml)


# burrow

Golang code for exploring some network centrality ideas. Still very much a WIP.

### Important dependencies:
* [Gonum](https://github.com/gonum/gonum): burrow is being built according to gonum Graph interface specs to provide access to the full range of Gonum graph algorithms.
* [Gomega](https://github.com/onsi/gomega): Matchers for testing.
* [Ginkgo](https://github.com/onsi/ginkgo): BDD testing framework for Go.


### Testing

Burrow tests are written using [Ginkgo](https://onsi.github.io/ginkgo), which runs on top of Go's native testing framework. To execute, run `ginkgo -r` or `go test ./...`.
