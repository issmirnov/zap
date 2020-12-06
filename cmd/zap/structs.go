package zap

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Jeffail/gabs/v2"
)

type Context struct {
	// Json container with path configs
	Config *gabs.Container

	// Enables safe hot reloading of Config.
	ConfigMtx sync.Mutex
}

type CtxWrapper struct {
	*Context
	H func(*Context, http.ResponseWriter, *http.Request) (int, error)
}

func (cw CtxWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := cw.H(cw.Context, w, r) // this runs the actual handler, defined in struct.
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
