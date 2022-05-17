package zap

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"encoding/json"

	"github.com/Jeffail/gabs/v2"
)

// IndexHandler handles all the non status expansions.
func IndexHandler(ctx *Context, w http.ResponseWriter, r *http.Request) (int, error) {
	var host string
	if r.Header.Get("X-Forwarded-Host") != "" {
		host = r.Header.Get("X-Forwarded-Host")
	} else {
		host = r.Host
	}

	var hostConfig *gabs.Container
	var ok bool

	// Check if host present in Config.
	children := ctx.Config.ChildrenMap()
	if hostConfig, ok = children[host]; !ok {
		return 404, fmt.Errorf("Shortcut '%s' not found in Config.", host)
	}

	tokens := tokenize(host + r.URL.Path)

	// Set up handles on token and Config. We might need to skip ahead if there's a custom schema set.
	tokensStart := tokens.Front()
	conf := ctx.Config

	var path bytes.Buffer
	if s := hostConfig.Path(sslKey).Data(); s != nil && s.(bool) {
		path.WriteString(httpPrefix)
	} else if s := hostConfig.Path(schemaKey).Data(); s != nil && s.(string) != "" {
		path.WriteString(hostConfig.Path(schemaKey).Data().(string) + ":/")
		// move one token ahead to parse expansions correctly.
		conf = conf.ChildrenMap()[tokensStart.Value.(string)]
		tokensStart = tokensStart.Next()
	} else {
		// Default to regular https prefix.
		path.WriteString(httpsPrefix)
	}

	ExpandPath(conf, tokensStart, &path)

	// send result
	http.Redirect(w, r, path.String(), http.StatusFound)

	return 302, nil
}

// HealthHandler responds to /healthz request.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `OK`)
}

// VarsHandler responds to /varz request and prints Config.
func VarsHandler(c *Context, w http.ResponseWriter, r *http.Request) (int, error) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, jsonPrettyPrint(c.Config.String()))
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
