package main

import (
    "bytes"
    "io"
    "log"
    "net/http"
    "strings"

    "github.com/Jeffail/gabs"
)

type context struct {
    // Json container with path configs
    config *gabs.Container
    // String to append to path for trailing slashes
    index string
}

type ctxWrapper struct {
    *context
    H func(*context, http.ResponseWriter, *http.Request) (int, error)
}

func (cw ctxWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    status, err := cw.H(cw.context, w, r) // this runs the actual handler, defined in struct.
    if err != nil {
        log.Printf("HTTP %d: %q", status, err)
        switch status {
        case http.StatusInternalServerError:
            http.Error(w, http.StatusText(status), status)
            // TODO - add bad request?
        default:
            http.Error(w, http.StatusText(status), status)
        }
    }
}

func IndexHandler(a *context, w http.ResponseWriter, r *http.Request) (int, error) {
    var host string
    if r.Header.Get("X-Forwarded-Host") != "" {
        host = r.Header.Get("X-Forwarded-Host")
    } else {
        host = r.Host
    }

    // massage the path, appending string to trailing slash if needed
    if strings.HasSuffix(r.URL.Path, "/") {
        r.URL.Path += a.index // TODO make this dependent on TLD
    }

    tokens := tokenize(host + r.URL.Path)
    var res bytes.Buffer
    res.WriteString("https:/") // second slash appended in expand() call
    expand(a.config,
        tokens.Front(), &res)

    // send result
    http.Redirect(w, r, res.String(), http.StatusFound)

    return 302, nil // not really needed?
}

// HealthHandler responds to /healthz request.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    io.WriteString(w, `OK`)
}

// StatusHandler responds to /statusz request and prints config.
func VarsHandler(c *context, w http.ResponseWriter, r *http.Request) (int, error) {
    w.WriteHeader(http.StatusOK)
    io.WriteString(w, "Config: "+c.config.String())
    return 200, nil
}
