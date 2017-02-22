# Zap (WIP)

A simple recursive URL expander.

Note: As of right now, this repo is still being set up. This message will be removed when everything is ready.

## Overview

ZAP is a simple go app that sends 302 redirects. That's it. It was written in just a few hours in between random household tasks. It helps people be more efficient by providing simple shortcuts for common pages.


## Configuration

### gotchas

- `n` and `y` are reserved for canonical yaml booleans. Use quotes in config to avoid issues.

## Benchmarks

TODO: add `ab` results, as well as `go test -bench .`

---

## Installation

### Quick Install

If you just want to hack around:

1. `go install github.com/issmirnov/zap`
2. `tmux new-session -t zap`
3. `$GOPATH/bin/zap`
4. `curl -I -L -H 'Host: g' localhost:8927/i`

If you want to actually install this properly, read on.

### Step 1: Set up Zap to run as a service



Nginx

```
server {
    listen 80;
    server_name s sp sg sf sd;
    error_log /var/log/nginx/error.log;
    access_log /var/log/nginx/access.log;

    location / {
        proxy_http_version 1.1;
        proxy_pass http://localhost:8927;
        proxy_set_header X-Forwarded-Host $host;
        #proxy_set_header Host $host;
    }
}
```



### Step 2: Configure DNS

For a lone machine, `/etc/hosts` is easiest. Add one entry per top level shortcut. For the sample config that would be `127.0.0.1 e` and `127.0.0.1 g`. The alternative is `dnsmasqd`, in which case you want to set up dnsmasqd to point to local server.

### Step 3: Tweak the config file and launch the service

- `brew services start zap`
- `systemctl start zap`

---


## Contributing

Standard GitHub workflow - fork and submit a PR.

### Handy Commands

- `go run main.go config.go text.go web.go` to run locally
- `curl -I -L -H 'Host: g' localhost:8927/i/dotfiles` - to test locally e2e
- `goconvey` - launches web UI for go tests.
- `go test` runs CLI tests.


## Tasks Roadmap

A short list of upcoming features and fixes, sorted by deadline.

- GoReleaser for release automation
- Homebrew install role
- Systemd service script
- configurable index page. so 'start' or 'index.html', set in top level domain config.
- queries: so `s/dns` -> `smirnov.wiki/start?do=search&id=dns`
- Travis CI
- coverage and go health report badges
- ansible role to facilitate installation
- add check for dual "expand" keys in config



## License

`ZAP` is licensed with MIT

## Contributors

- [Ivan Smirnov](http://ivansmirnov.name)
