package hkit

import (
	"net/http"
)

type responseWriter struct {
	headerFunc      func() http.Header
	writeFunc       func([]byte) (int, error)
	writeHeaderFunc func(int)
}

func (w *responseWriter) Header() http.Header {
	return w.headerFunc()
}

func (w *responseWriter) Write(bytes []byte) (int, error) {
	return w.writeFunc(bytes)
}

func (w *responseWriter) WriteHeader(header int) {
	w.writeHeaderFunc(header)
}
