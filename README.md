# Zap (WIP)

[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)

A simple recursive URL expander.

**Note: As of right now, this repo is still being set up. This message will be removed when everything is ready.**

## Overview

ZAP is a simple go app that sends 302 redirects. That's it. It was written in just a few hours in between random household tasks. It helps people be more efficient by providing simple shortcuts for common pages.

It was fun to build because it's super lightweight and insanely fast. Some sample benchmarks:

```
# Trial 1: localhost
$ ab -n 10000 -c 100 http://localhost:8927/
Requests per second:    39888.31 [#/sec] (mean)
Time per request:       2.507 [ms] (mean)

# Trial 2: Hitting server on LAN over gigabit, zap behind nginx proxy
$ ab -n 10000 -c 100 http://server/z
Requests per second:    12671.57 [#/sec] (mean)
Time per request:       7.892 [ms] (mean)

# Go benchmarks
$ go test -bench=.
BenchmarkIndexHandler-8   	 1000000	      1679 ns/op
```

As you can see, even behind an nginx proxy, running on a server that I
mercilessly abuse we still get a respectable ~13k QPS.



## Installation

### Step 0: Quick Install

If you just want to hack around:

1. `go install github.com/issmirnov/zap`
2. `tmux new-session -t zap`
3. `$GOPATH/bin/zap`
4. `curl -I -L -H 'Host: g' localhost:8927/z` - should return 302 to github.com/issmirnov.zap

If you want to actually install this properly, read on.


### Step 1: Set up Zap to run as a service

Zap takes several command line flags:

- `-config` - path to config file. Default is `./c.yml`
- `-port` - port to bind to. Use either 80 or 8927

#### OSX (brew)

1. `brew install issmirnov/apps/zap`
2. `sudo brew services start zap`

If you already have port 80 in use, you can run zap behind a reverse proxy.

1. Change the port in the zap plist config: `sed -i '' 's/8927/80/'  /usr/local/Cellar/zap/*/homebrew.mxcl.zap.plist`
2. Start zap as user service: `brew services start zap`
3. Configure your web server to act as a reverse proxy. For nginx, this will suffice:

```
# File: /usr/local/etc/nginx/servers/zap.conf
server {
    listen 80; # Keep as 80, make sure nginx listens on 80 too
    server_name e g; # Put your shortcuts here
location / {
        proxy_http_version 1.1;
        proxy_pass http://localhost:8927;
        proxy_set_header X-Forwarded-Host $host;
        #proxy_set_header Host $host;
    }
}
```

4. Restart your webserver and test the result: `curl -I -L -H 'Host: g' localhost:8927/z`

#### Ubuntu

TODO: Provide systemd script, add hints about nginx.


### Step 2: Configure DNS

If you are running `zap` locally, you need to edit `/etc/hosts` and add each top level entry. For the sample config, this would be `127.0.0.1 e` and `127.0.0.1 g`. Adjust accordingly.

For the advanced users, I suggest running `dnsmasqd` and add DNS entries for all TLDs, so that all of your clients will automatically point to the server.

### Step 3: Tweak the config file and launch the service

The config file is located at `/usr/local/etc/zap/` on OSX, and `/path/tbd` on ubuntu.

Open up `c.yml` and update the mappings you would like. You can nest arbitrarily deep. Expansions work on strings and ints.

Important gotcha: yaml has [reserved types](http://yaml.org/type/bool.html) and thus `n`, `y`, `no` and the like need to be quoted. See the sample config.

Once you're done editing the file (making sure to keep the DNS entries in sync) restart the service and test it out.

- OSX: `sudo brew services restart zap` or `brew services restart zap`
- Ubuntu: `systemctl restart zap`

## Contributing

Patches are welcome! Please use the standard GitHub workflow - fork this repo and submit a PR. I'll usually get to it within a few days.

Handy commands for local dev:

- `go run main.go config.go text.go web.go` to run locally
- `curl -I -L -H 'Host: g' localhost:8927/z` - to test locally e2e
- `goconvey` - launches web UI for go tests.
- `go test` runs CLI tests.


## Tasks Roadmap

A short list of upcoming features and fixes, sorted by deadline.

- GoReleaser for release automation - DONE
- Homebrew install script - DONE
- Systemd service script
- configurable index page. so 'start' or 'index.html', set in top level domain config.
- queries: so `s/dns` -> `smirnov.wiki/start?do=search&id=dns`
- Travis CI
- coverage and go health report badges
- ansible role to facilitate installation
- add check for dual "expand" keys in config


## Contributors

- [Ivan Smirnov](http://ivansmirnov.name)
