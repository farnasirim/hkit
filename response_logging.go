package hkit

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

type responseLogger struct {
	next       http.HandlerFunc
	bodyLogger *responseBodyLogger
}

func newResponseLogger(handlerFunc http.HandlerFunc) *responseLogger {
	bodyLogger := newResponseBodyLogger(handlerFunc)

	logger := &responseLogger{
		next:       bodyLogger.ServeHTTP,
		bodyLogger: bodyLogger,
	}

	return logger
}

func (l *responseLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.next(w, r)
}

func (l *responseLogger) SetWriter(writer io.Writer) *responseLogger {
	l.bodyLogger.SetWriter(writer)
	return l
}

type responseBodyLogger struct {
	next   http.HandlerFunc
	writer io.Writer
}

func newResponseBodyLogger(handlerFunc http.HandlerFunc) *responseBodyLogger {
	logger := &responseBodyLogger{
		next:   handlerFunc,
		writer: os.Stdout,
	}

	return logger
}

func (l *responseBodyLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpWriter := bytes.NewBuffer(nil)
	bodyLoggerResponseWriter := &responseWriter{
		headerFunc: func() http.Header {
			// We don't care about the header here
			return make(http.Header)
		},
		writeFunc: func(bytes []byte) (int, error) {
			return tmpWriter.Write(bytes)
		},
		writeHeaderFunc: func(int) {
			// No need to do anything
		},
	}
	writer := NewLazyMultiWriter(w, bodyLoggerResponseWriter)
	l.next(writer, r)
	l.writer.Write(tmpWriter.Bytes())
	l.writer.Write([]byte("\n\n"))
}

func (l *responseBodyLogger) SetWriter(writer io.Writer) *responseBodyLogger {
	l.writer = writer
	return l
}
