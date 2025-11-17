package zap

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Jeffail/gabs/v2"
)

type Context struct {
	// Config is a Json container with path configs
	Config *gabs.Container

	// ConfigMtx Enables safe hot reloading of Config.
	ConfigMtx sync.Mutex

	// Advertise IP, used in /etc/hosts in case bind address differs.
	Advertise string
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
			// Provide more helpful error messages for internal server errors
			errorMsg := fmt.Sprintf("Internal Server Error: %s", err.Error())
			http.Error(w, errorMsg, status)
		case http.StatusNotFound:
			// Provide helpful message for 404 errors
			errorMsg := fmt.Sprintf("Shortcut not found: %s", err.Error())
			http.Error(w, errorMsg, status)
		case http.StatusBadRequest:
			// Provide helpful message for 400 errors
			errorMsg := fmt.Sprintf("Bad Request: %s", err.Error())
			http.Error(w, errorMsg, status)
		default:
			// For other status codes, use the error message directly
			http.Error(w, err.Error(), status)
		}
	}
}
