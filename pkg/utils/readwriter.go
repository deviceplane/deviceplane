package utils

import (	
	"io"
	"net/http"
)

type ResponseWriter struct {
	HasWrittenStatus bool

	Headers http.Header
	Writer  io.Writer
	Status  int
}

func (w *ResponseWriter) Write(b []byte) (n int, err error) {
	if !w.HasWrittenStatus {
		w.WriteHeader(http.StatusOK)
	}

	return w.Writer.Write(b)
}

func (w *ResponseWriter) Header() http.Header {
	return w.Headers
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.Status = code
	w.HasWrittenStatus = true
}
