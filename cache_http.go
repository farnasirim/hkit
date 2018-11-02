package hkit

import (
	"errors"
	"net/http"
)

type HTTPCache struct {
	cacheService CacheService
}

var (
	errHeaderAlreadyLocked = errors.New("Header has already been locked")
	ErrAlreadyWritten      = errors.New("Response has already been written")
)

// HTTPResponseInterceptor will capture actions on an underlaying
// ResponseWriter. IMPORTANT NOTE: will only capture header state until
// the first call to WriteHeader or Write. Therefore no support for
// trailers etc.
// Behavior may diverge from the implementors in the "net/http", so no
// crazily protective behaviors should be expected from this wrapper.
// TL;DR: you probably have your own cache wrapper if you are writing
// e.g. trailer headers all over your handlers.
type HTTPResponseInterceptor struct {
	tempHeader   http.Header
	lockedHeader http.Header
	writtenBytes []byte
	writtenCode  int
}

func NewHTTPResponseInterceptor(w http.ResponseWriter) *HTTPResponseInterceptor {
	interceptor := &HTTPResponseInterceptor{
		tempHeader:   cloneHeader(w.Header()),
		lockedHeader: nil,
		writtenBytes: nil,
		writtenCode:  -1,
	}
	return interceptor
}

func (w *HTTPResponseInterceptor) Write(bytes []byte) (int, error) {
	w.lockHeader()
	if w.writtenBytes != nil {
		return 0, ErrAlreadyWritten
	}
	w.writtenBytes = bytes
	return w.Write(bytes)
}

func (w *HTTPResponseInterceptor) WriteHeader(code int) {
	if err := w.lockHeader(); err != nil {
		return
	}
	w.writtenCode = code
	w.WriteHeader(code)
}

func (w *HTTPResponseInterceptor) Header() http.Header {
	return w.tempHeader
}

func (w *HTTPResponseInterceptor) lockHeader() error {
	if w.lockedHeader != nil {
		return errHeaderAlreadyLocked
	}
	w.lockedHeader = cloneHeader(w.tempHeader)
	overwriteHeader(w.Header(), w.lockedHeader)
	return nil
}

func (w *HTTPResponseInterceptor) Replay(writer http.ResponseWriter) {
	if w.lockedHeader != nil {
		overwriteHeader(writer.Header(), w.lockedHeader)
	}
	if w.writtenCode != -1 {
		writer.WriteHeader(w.writtenCode)
	}
	if w.writtenBytes != nil {
		writer.Write(w.writtenBytes)
	}
}

func cloneHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}

func overwriteHeader(lhs, rhs http.Header) {
	keys := make([]string, 0)
	for k, _ := range lhs {
		keys = append(keys, k)
	}
	for _, k := range keys {
		lhs.Del(k)
	}

	for k, v := range rhs {
		vv := make([]string, len(v))
		copy(vv, v)
		lhs[k] = vv
	}
}
