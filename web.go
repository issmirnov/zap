package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"encoding/json"
	"github.com/Jeffail/gabs"
)

type context struct {
	// Json container with path configs
	config *gabs.Container

	// Enables safe hot reloading of config.
	configMtx sync.Mutex
}

type ctxWrapper struct {
	*context
	H func(*context, http.ResponseWriter, *http.Request) (int, error)
}

func (cw ctxWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := cw.H(cw.context, w, r) // this runs the actual handler, defined in struct.
	if err != nil {
		switch status {
		case http.StatusInternalServerError:
			http.Error(w, fmt.Sprintf("HTTP %d: %q", status, err), status)
			// TODO - add bad request?
		default:
			http.Error(w, err.Error(), status)
		}
	}
}

// IndexHandler handles all the non status expansions.
func IndexHandler(a *context, w http.ResponseWriter, r *http.Request) (int, error) {
	var host string
	if r.Header.Get("X-Forwarded-Host") != "" {
		host = r.Header.Get("X-Forwarded-Host")
	} else {
		host = r.Host
	}

	var hostConfig *gabs.Container
	var ok bool

	// Check if host present in config.
	children, _ := a.config.ChildrenMap()
	if hostConfig, ok = children[host]; !ok {
		return 404, fmt.Errorf("Shortcut '%s' not found in config.", host)
	}

	tokens := tokenize(host + r.URL.Path)
	var path bytes.Buffer
	if s := hostConfig.Path(sslKey).Data(); s != nil && s.(bool) {
		path.WriteString(httpPrefix)
	} else {
		path.WriteString(httpsPrefix)
	}

	expandPath(a.config, tokens.Front(), &path)

	// send result
	http.Redirect(w, r, path.String(), http.StatusFound)

	return 302, nil
}

// HealthHandler responds to /healthz request.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `OK`)
}

// VarsHandler responds to /varz request and prints config.
func VarsHandler(c *context, w http.ResponseWriter, r *http.Request) (int, error) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, jsonPrettyPrint(c.config.String()))
	return 200, nil
}

// https://stackoverflow.com/a/36544455/5117259
func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	out.WriteString("\n")
	return out.String()
}
