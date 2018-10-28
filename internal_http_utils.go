package hkit

import (
	"errors"
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

var (
	// ErrDifferentBytesWritten is returned when we write different number of bytes
	// to each of the http.ResponseWriter s passed to the LazyMultiWriter
	ErrDifferentBytesWritten = errors.New(`Different number of bytes written while\
writing multiple headers`)
)

// LazyMultiWriter is an `http.ResponseWriter` which simulates its opeartions
// under the hood on multiple different `http.ResponseWriter`s
type LazyMultiWriter struct {
	header  http.Header
	writers []http.ResponseWriter
}

// Header returns the internal http.Header object in the `LazyMultiHeader`
// This header will be written to the underlaying `http.ResponseWriter`s
// Upon calling `Finish` on the `LazyMultiHeader`
func (w *LazyMultiWriter) Header() http.Header {
	return w.header
}

func (w *LazyMultiWriter) Write(bytes []byte) (int, error) {
	bytesWritten := -1
	var err error
	for _, writer := range w.writers {
		ret, _ := writer.Write(bytes)
		if bytesWritten != -1 && ret != bytesWritten {
			err = ErrDifferentBytesWritten
			break
		}
		bytesWritten = ret
	}
	return bytesWritten, err
}

// WriteHeader calls WriteHeader of every http.ResponseWriter that has been
// passed to the LazyMultiWriter
func (w *LazyMultiWriter) WriteHeader(header int) {
	for _, writer := range w.writers {
		writer.WriteHeader(header)
	}
}

func overwriteHeader(lhs http.Header, rhs http.Header) {
	for key := range lhs {
		lhs.Del(key)
	}

	for k, vv := range rhs {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		lhs[k] = vv2
	}
}

// Finish writes the accumulated header to each of the http.ResponseWriters
// YOU HAVE TO MANUALLY CALL FINISH SOMETIME. This is done in such an ugly
// to avoid reimplementing the whole http.Header type
func (w *LazyMultiWriter) Finish() {
	for _, writer := range w.writers {
		overwriteHeader(writer.Header(), w.header)
	}
}

// NewLazyMultiWriter creates a new LazyMultiWriter which is used for writing
func NewLazyMultiWriter(writers ...http.ResponseWriter) *LazyMultiWriter {
	return &LazyMultiWriter{
		header:  make(http.Header),
		writers: writers,
	}
}
