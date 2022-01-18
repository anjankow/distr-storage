package server

import (
	"net/http"
)

const (
	contentType = "application/json"
)

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("all good here"))
}
