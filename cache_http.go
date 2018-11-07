package hkit

import (
	"bytes"
	"encoding/gob"
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

// HTTPResponseReplayer will capture actions on an underlaying
// ResponseWriter, or replay already captured actions.
// IMPORTANT NOTE: will only capture header state until the first call to
// WriteHeader or Write. Therefore no support for trailers etc.
// Should be used once. Instantiate a new one for each time you're proxying
// the write process. E.g. instantiate inside your middleware function per req.
// Behavior may diverge from the implementors in the "net/http", so no
// crazily protective behaviors should be expected from this wrapper.
// TL;DR: you probably have your own cache wrapper if you are writing
// e.g. trailer headers all over your handlers.
type HTTPResponseReplayer struct {
	tempHeader   http.Header
	next         http.ResponseWriter
	lockedHeader http.Header
	writtenBytes []byte
	writtenCode  int
}

func NewHTTPResponseReplayer(
	w http.ResponseWriter) *HTTPResponseReplayer {

	replayer := &HTTPResponseReplayer{
		tempHeader:   cloneHeader(w.Header()),
		next:         w,
		lockedHeader: nil,
		writtenBytes: make([]byte, 0),
		writtenCode:  -1,
	}
	return replayer
}

func (w *HTTPResponseReplayer) Serialize(buffer []byte) (int, error) {
	return 0, nil
}

func ReplayerFromSerialized(savedState []byte) *HTTPResponseReplayer {
	replayer := NewHTTPResponseReplayer(nil)

	return replayer
}

func (w *HTTPResponseReplayer) Write(bytes []byte) (int, error) {
	w.lockHeader()
	w.writtenBytes = append(w.writtenBytes, bytes...)
	return w.next.Write(bytes)
}

func (w *HTTPResponseReplayer) WriteHeader(code int) {
	if err := w.lockHeader(); err != nil {
		return
	}
	w.writtenCode = code
	w.next.WriteHeader(code)
}

func (w *HTTPResponseReplayer) Header() http.Header {
	return w.tempHeader
}

func (w *HTTPResponseReplayer) lockHeader() error {
	if w.lockedHeader != nil {
		return errHeaderAlreadyLocked
	}
	w.lockedHeader = cloneHeader(w.tempHeader)
	overwriteHeader(w.next.Header(), w.lockedHeader)
	return nil
}

func (w *HTTPResponseReplayer) Replay(writer http.ResponseWriter) {
	if w.lockedHeader != nil {
		overwriteHeader(writer.Header(), w.lockedHeader)
	}
	if w.writtenCode != -1 {
		writer.WriteHeader(w.writtenCode)
	}
	writer.Write(w.writtenBytes)
}

func (w *HTTPResponseReplayer) IsHTTPStatusOK() bool {
	return w.lockedHeader != nil && (w.writtenCode == -1 ||
		w.writtenCode == http.StatusOK)
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

type responseReplayerPayload struct {
	TempHeader   http.Header
	LockedHeader http.Header
	WrittenBytes []byte
	WrittenCode  int
}

func newResponseReplayerPayload(tempHeader, lockedHeader http.Header,
	writtenBytes []byte, writtenCode int) *responseReplayerPayload {
	return &responseReplayerPayload{
		TempHeader:   tempHeader,
		LockedHeader: lockedHeader,
		WrittenBytes: writtenBytes,
		WrittenCode:  writtenCode,
	}
}

func marshalReplayer(r *HTTPResponseReplayer) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0))
	gob.NewEncoder(buffer).Encode(newResponseReplayerPayload(
		r.tempHeader,
		r.lockedHeader,
		r.writtenBytes,
		r.writtenCode,
	))
	return buffer.Bytes()
}

func unmarshalReplayer(savedState []byte) (*HTTPResponseReplayer, error) {
	decoded := &responseReplayerPayload{}
	err := gob.NewDecoder(bytes.NewBuffer(savedState)).Decode(decoded)
	if err != nil {
		return nil, err
	}
	return &HTTPResponseReplayer{
		tempHeader:   decoded.TempHeader,
		lockedHeader: decoded.LockedHeader,
		writtenBytes: decoded.WrittenBytes,
		writtenCode:  decoded.WrittenCode,
	}, nil
}
