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
	// Check if context and configuration are valid
	if ctx == nil || ctx.Config == nil {
		return http.StatusInternalServerError, fmt.Errorf("server configuration is invalid or not loaded")
	}

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
		return http.StatusNotFound, fmt.Errorf("shortcut '%s' not found in config", host)
	}

	tokens := tokenize(host + r.URL.Path)

	// Set up handles on token and Config. We might need to skip ahead if there's a custom schema set.
	tokensStart := tokens.Front()
	conf := ctx.Config

	var path bytes.Buffer
	if s := hostConfig.Path(sslKey).Data(); s != nil && s.(bool) {
		path.WriteString(httpPrefix)
	} else if s := hostConfig.Path(schemaKey).Data(); s != nil && s.(string) != "" {
		schema := hostConfig.Path(schemaKey).Data().(string)
		if schema == "" {
			return http.StatusInternalServerError, fmt.Errorf("invalid schema configuration for host '%s'", host)
		}
		path.WriteString(schema + ":/")
		// move one token ahead to parse expansions correctly.
		conf = conf.ChildrenMap()[tokensStart.Value.(string)]
		tokensStart = tokensStart.Next()
	} else {
		// Default to regular https prefix.
		path.WriteString(httpsPrefix)
	}

	// Validate that we have a valid configuration before expanding
	if conf == nil {
		return http.StatusInternalServerError, fmt.Errorf("invalid configuration structure for host '%s'", host)
	}

	if err := ExpandPath(conf, tokensStart, &path); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to expand path for host '%s': %w", host, err)
	}

	// Validate that we generated a valid path
	if path.Len() == 0 {
		return http.StatusInternalServerError, fmt.Errorf("failed to generate redirect path for host '%s'", host)
	}

	// send result
	http.Redirect(w, r, path.String(), http.StatusFound)

	return http.StatusFound, nil
}

// HealthHandler responds to /healthz request.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, `OK`)
	if err != nil {
		// Log the error for debugging but still return a 500
		http.Error(w, "Internal Server Error: Health check failed", http.StatusInternalServerError)
	}
}

// VarsHandler responds to /varz request and prints Config.
func VarsHandler(c *Context, w http.ResponseWriter, r *http.Request) (int, error) {
	// Validate that we have a valid configuration
	if c.Config == nil {
		return http.StatusInternalServerError, fmt.Errorf("configuration not loaded or invalid")
	}

	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, jsonPrettyPrint(c.Config.String()))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to write response: %w", err)
	}
	return http.StatusOK, nil
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
