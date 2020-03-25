package service

import (
	"encoding/base64"
	"net/http"
	"strconv"
)

func withPort(w http.ResponseWriter, r *http.Request, f func(port int)) {
	query := r.URL.Query()

	portRaw := query.Get("port")
	port, err := strconv.Atoi(portRaw)
	if err != nil {
		http.Error(w, "invalid port", 400)
		return
	}

	f(port)
}

func withPath(w http.ResponseWriter, r *http.Request, f func(path string)) {
	query := r.URL.Query()

	pathRaw := query.Get("path")
	path, err := base64.RawURLEncoding.DecodeString(pathRaw)
	if err != nil {
		http.Error(w, "invalid base64 encoded path", 400)
		return
	}

	f(string(path))
}
