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

ZAP is a simple go app that sends 302 redirects. It's insanely fast, maxing out at over 150k qps. It helps people be more efficient by providing simple shortcuts for common pages.

It can help save keystores on any level of the URL. In the example above, the user types `gh/z` and zap expands `gh` into `github.com` and `z` into `issmirnov/zap`. There is no limit to how deep you can go. Zap can be useful for shortening common paths. If your or your company has many projects at `company.com/long/and/annoying/path/name_here` zap can turn this into `c/name_here`, or `c/p/name_here` - it's all in your hands. 

Zap runs as an HTTP service, and can live on the standard web ports or behind a proxy. It features hot reloading of the config, super low memory footprint and amazing durability under heavy loads.

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

When you add a new shortcut, you need to indicate to your web browser that it's not a search term. You can do this by typing it in once with just a slash. For example, if you add a shortcut `g/z` -> github.com/issmirnov/zap, if you try `g/z` right away you will get taken to the search page. Instead, try `g/` once, and then `g/z`. This initial step only needs to be taken once per new shortcut.

#### Examples

You can configure your `c.yml` file endlessly. Here are some examples to get inspire your creativity:

```
g:
  expand: github.com
  z:
    expand: github.com/issmirnov/zap
f:
  expand: facebook.com
  php:
    expand: groups/2204685680/
r:
  expand: reddit.com/r
```

With this config, you can use the following queries:

  - `g/z` -> github.com/issmirnov/zap
  - `f/zuck` -> facebook.com/zuck
  - `f/php` -> facebook.com/groups/2204685680/
  - `r/catsstandingup` -> reddit.com/r/catsstandingup

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

Benchmarked with [wrk2](https://github.com/giltene/wrk2) on Ubuntu 16.04 using an i5 4590 CPU.
```
# Maxing out QPS.
$ wrk -t2 -c10 -d30s -R500000  http://127.0.0.1:8989/h
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     3.40s     1.96s    6.79s    57.82%
    Req/Sec    80.22k   320.00    80.57k    50.00%
Requests/sec: 161077.54

# Getting max users while longest request under 15ms
$ ./wrk -t2 -c10 -d20s -R120000 http://127.0.0.1:8989/h
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.15ms    0.93ms  14.38ms   81.79%
    Req/Sec    63.24k     7.18k  110.89k    76.92%
Requests/sec: 119932.12

```
As you can see, zap peaks at around ~160k qps, and can sustain ~120k qps with an average response under 15ms.

Note: The config used was:

```
e:
  expand: example.com
  a:
    expand: apples
  b:
    expand: bananas
g:
  expand: github.com
  d:
    expand: issmirnov/dotfiles
  s:
    query: search?q=
  z:
    expand: issmirnov/zap
'127.0.0.1:8989':
  expand: '127.0.0.1:8989'
  h:
    expand: healthz
```

## Contributing

Patches are welcome! Please use the standard GitHub workflow - fork this repo and submit a PR. I'll usually get to it within a few days.

Handy commands for local dev:

- `go run main.go config.go text.go web.go` to run locally
- `curl -I -L -H 'Host: g' localhost:8927/z` - to test locally e2e
- `goconvey` - launches web UI for go tests.
- `go test` runs CLI tests.

## Contributors

- [Ivan Smirnov](http://ivansmirnov.name)
