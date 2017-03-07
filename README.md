# Zap

[![Travis](https://travis-ci.org/issmirnov/zap.svg?branch=master)](https://travis-ci.org/issmirnov/zap)
[![Release](https://img.shields.io/github/release/issmirnov/zap.svg?style=flat-square)](https://github.com/issmirnov/zap/releases/latest)
![Total Downloads](https://img.shields.io/github/downloads/issmirnov/zap/total.svg)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/issmirnov/zap?style=flat-square)](https://goreportcard.com/report/github.com/issmirnov/zap)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)

A simple URL redirector. Allows you to define shortcuts to pages. I find it faster than traditional bookmarks.

![zap demo gif](zap_demo.gif)

## Overview

ZAP is a simple go app that sends 302 redirects. That's it. It was written in just a few hours in between random household tasks. It helps people be more efficient by providing simple shortcuts for common pages.

## Installation

### Ansible

If you know how to use ansible, head over to the [ansible galaxy](https://galaxy.ansible.com/issmirnov/zap/) and install this role as `issmirnov.zap`. I've done the heavy lifting for you. If you want to do this by hand, read on.

### OSX: brew install

1. `brew install issmirnov/apps/zap`
2. `sudo brew services start zap`
3. Add shortcuts to `/usr/local/etc/zap/c.yml`
3. Enjoy!

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

Restart your webserver and test the result: `curl -I -L -H 'Host: g' localhost/z`



### Ubuntu

Note: This section applies to systemd systems only. If you are running ubuntu 14.10 or below, you'll have to use [upstart](https://www.digitalocean.com/community/tutorials/the-upstart-event-system-what-it-is-and-how-to-use-it).

1. create the zap user: `sudo adduser --system --no-create-home --group zap`

2. Make a service definition at `/etc/systemd/system/zap.service`. If you are running zap behind a web server, use the following config:

```
[Unit]
Description=Zap (URL text expander)
After=syslog.target
After=network.target

[Service]
Type=simple
User=zap
Group=zap
WorkingDirectory=/etc/zap
ExecStart=/usr/local/bin/zap -port 8927 -config c.yml
Restart=always
RestartSec=2s

[Install]
WantedBy=multi-user.target
```

If you are running standalone:
```
[Unit]
Description=Zap (URL text expander)
After=syslog.target
After=network.target

[Service]
Type=simple
User=root
Group=zap
WorkingDirectory=/etc/zap
ExecStart=/usr/local/bin/zap -port 80 -config c.yml
Restart=always
RestartSec=2s

[Install]
WantedBy=multi-user.target
```

You'll notice the difference is that we have to run as `root` in order to bind to port 80. *If you know of a way to launch a go app under setuid and then drop privileges, please send a PR*

3. Start your new service: `sudo systemctl start zap` and make sure it's running: `sudo systemctl status zap`


### Configuration

The config file is located at `/usr/local/etc/zap/c.yml` on OSX. For ubuntu, you will have to create `/etc/zap/c.yml` by hand.

Open up `c.yml` and update the mappings you would like. You can nest arbitrarily deep. Expansions work on strings and ints. Notice that we have two keywords available: `expand` and `query`. The `query` term acts almost like the `expand` option, but drops the separating slash between query expansion and search term (`example.com?q=foo` instead of `example.com?q=/foo`)

Important gotcha: yaml has [reserved types](http://yaml.org/type/bool.html) and thus `n`, `y`, `no` and the like need to be quoted. See the sample config.

Zap supports hot reloading, so simply save the file when you are done and test out your new shortcut. Note: If the shortcut does not work, make sure your YAML is correct and that zap is not printing any errors. You can test this by stopping zap and starting it manually - it should print any issues to stdout.

For the advanced users: You might have to reload your webserver and `dnsmasq`, depending on your setup.

### Additional Information


#### Zap flags:

- `-config` - path to config file. Default is `./c.yml`
- `-port` - port to bind to. Default is 8927. Use 80 in standalone mode.
- `-host` - default is 127.0.0.1. Use 0.0.0.0 for a public server.

### DNS management via /etc/hosts

Zap will attempt to keep the `/etc/hosts` file in sync with the configuration specified. This is assumed to be a reasonable default. If you wish to disable this behavior, run zap under a user that does not have write permissions to that file.

As long as you don't touch the delimiters used by zap (`### Zap Shortcuts :start ##` and `### Zap Shortcuts :end ##`) you can edit the hosts file as you wish. If those delimiters are missing, zap will append them to the file. You shouldn't have problems if you manage your hosts file in a reasonable manner.

For the advanced users running zap on a server on an internal network, I suggest looking into `dnsmasqd` - this will allow all your clients to utilize these shortcuts globally.


## Benchmarks

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

When we run zap locally we easily hit almost 40k qps. Even when behind an nginx proxy, running on a server that I
mercilessly abuse we still get a respectable ~13k QPS.

## Contributing

Patches are welcome! Please use the standard GitHub workflow - fork this repo and submit a PR. I'll usually get to it within a few days.

Handy commands for local dev:

- `go run main.go config.go text.go web.go` to run locally
- `curl -I -L -H 'Host: g' localhost:8927/z` - to test locally e2e
- `goconvey` - launches web UI for go tests.
- `go test` runs CLI tests.

## Contributors

- [Ivan Smirnov](http://ivansmirnov.name)
