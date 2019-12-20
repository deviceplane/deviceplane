package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

const (
	ProxiedFromDeviceHeader = "proxied-from-device"
)

var (
	errInvalidReferrer = errors.New("invalid referrer")
)

func JSONConvert(src, target interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &target)
}

func Respond(w http.ResponseWriter, ret interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ret)
}

func ProxyResponseFromDevice(w http.ResponseWriter, resp *http.Response) {
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.Header().Set(ProxiedFromDeviceHeader, "")

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

func ProxyResponse(w http.ResponseWriter, resp *http.Response) {
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

func WithReferrer(w http.ResponseWriter, r *http.Request, f func(referrer *url.URL)) {
	referrer, err := url.Parse(r.Referer())
	if err != nil {
		http.Error(w, errInvalidReferrer.Error(), http.StatusBadRequest)
		return
	}
	if referrer.Scheme != "http" && referrer.Scheme != "https" {
		http.Error(w, errInvalidReferrer.Error(), http.StatusBadRequest)
		return
	}
	f(referrer)
}
