## Contributing to Zap

Patches are welcome! Please use the standard GitHub workflow - fork this
repo and submit a PR. I'll usually get to it within a few days.

## Setting up dev environment

- Install [Goland](https://www.jetbrains.com/go/), [Atom](https://atom.io/),
or your favorite web editor with Golang support.

```
cd $GOPATH/src/github.com/issmirnov/zap
go build . # sanity check

# install test deps
go get github.com/smartystreets/goconvey
```

## Handy commands for local development:

- `go run main.go config.go text.go web.go` to run locally
- `curl -I -L -H 'Host: g' localhost:8927/z` - to test locally e2e
- `goconvey -excludedDirs dist` - launches web UI for go tests.
- `go test` runs CLI tests.


## Contributors

- [Ivan Smirnov](http://ivansmirnov.name)
- [Sergey Smirnov](https://smirnov.nyc/)
