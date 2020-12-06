## Contributing to Zap

Patches are welcome! Please use the standard GitHub workflow - fork this
repo and submit a PR. I'll usually get to it within a few days. If I miss your PR email, feel free to ping me directly at isgsmirnov@gmail.com after about a week.

## Setting up dev environment

- Install [Goland](https://www.jetbrains.com/go/), [Atom](https://atom.io/),
or your favorite web editor with Golang support.
- Note: This project relies on Go Modules, introduces in Go 1.11+.

```
git clone $your_fork zap
cd zap
go get
go build -o zap -v ./cmd/ # sanity check

# install test deps and run all tests
go test -short -v ./... -race -coverprofile=coverage.txt -covermode=atomic ./cmd 
./e2e.sh
```

## Handy commands for local development:

- `go build -o zap -v ./cmd/ && ./zap` to run locally
- `curl -I -L -H 'Host: g' localhost:8927/z` - to test locally e2e
- `goconvey -excludedDirs dist` - launches web UI for go tests.
- `./e2e.sh` runs CLI tests.


## Contributors

- [Ivan Smirnov](http://ivansmirnov.name)
- [Sergey Smirnov](https://smirnov.nyc/)
- [Chris Egerton](https://github.com/C0urante)