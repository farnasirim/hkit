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

// LazyMultiWriter is an `http.ResponseWriter` which simulates its opeartions
// under the hood on multiple different `http.ResponseWriter`s.
// WE DON'T CARE ABOUT ANY HEADER OTHER THAN THAT OF THE FIRST WRITER PASSED
type LazyMultiWriter struct {
	writers []http.ResponseWriter
}

// Header returns the header of the first writer passed
// zero writers passed? will panic. Sorry.
func (w *LazyMultiWriter) Header() http.Header {
	return w.writers[0].Header()
}

func (w *LazyMultiWriter) Write(bytes []byte) (int, error) {
	for _, writer := range w.writers {
		_, _ = writer.Write(bytes)
	}
	return len(bytes), nil
}

// WriteHeader calls WriteHeader of every http.ResponseWriter that has been
// passed to the LazyMultiWriter
func (w *LazyMultiWriter) WriteHeader(header int) {
	for _, writer := range w.writers {
		writer.WriteHeader(header)
	}
}

// NewLazyMultiWriter creates a new LazyMultiWriter which is used for writing
// ONLY THE `.Header()` OBJECT OF THE FIRST WRITER IS USED. PASS THE ORIGINAL
// WRITER AS THE FIRST ARGUMENT.
func NewLazyMultiWriter(writers ...http.ResponseWriter) *LazyMultiWriter {
	return &LazyMultiWriter{
		writers: writers,
	}
}
