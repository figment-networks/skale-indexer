// +build !profiling

package main

import (
	"net/http"
)

func attachProfiling(mux *http.ServeMux) {
	// noop
}
