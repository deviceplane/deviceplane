package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"unicode/utf8"
)

const (
	ProxiedFromDeviceHeader = "proxied-from-device"
)

var (
	errInvalidReferrer = errors.New("invalid referrer")
)

// From https://github.com/gorilla/websocket
func CheckSameOrAllowedOrigin(r *http.Request, validOrigins []url.URL) bool {
	originHeader := r.Header["Origin"]
	if len(originHeader) == 0 {
		return true
	}
	origin, err := url.Parse(originHeader[0])
	if err != nil {
		return false
	}

	if EqualASCIIFold(origin.Host, r.Host) {
		return true
	}
	for _, validOrigin := range validOrigins {
		if EqualASCIIFold(origin.Host, validOrigin.Host) {
			return true
		}
	}
	return false
}

// From https://github.com/gorilla/websocket
// EqualASCIIFold returns true if s is equal to t with ASCII case folding as
// defined in RFC 4790.
func EqualASCIIFold(s, t string) bool {
	for s != "" && t != "" {
		sr, size := utf8.DecodeRuneInString(s)
		s = s[size:]
		tr, size := utf8.DecodeRuneInString(t)
		t = t[size:]
		if sr == tr {
			continue
		}
		if 'A' <= sr && sr <= 'Z' {
			sr = sr + 'a' - 'A'
		}
		if 'A' <= tr && tr <= 'Z' {
			tr = tr + 'a' - 'A'
		}
		if sr != tr {
			return false
		}
	}
	return s == t
}

func InternalTags(projectID string) []string {
	return []string{
		"project:" + projectID,
	}
}

// Elliot Chance's github gist: https://gist.github.com/elliotchance/d419395aa776d632d897
func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}

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
