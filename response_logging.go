package hkit

import (
	"io"
	"net/http"
	"os"
)

type responseLogger struct {
	next         http.HandlerFunc
	bodyLogger   *responseBodyLogger
	headerLogger *responseHeaderLogger
}

func newResponseLogger(handlerFunc http.HandlerFunc) *responseLogger {
	bodyLogger := newResponseBodyLogger(handlerFunc)
	headerLogger := newResponseHeaderLogger(bodyLogger.ServeHTTP)

	logger := &responseLogger{
		next:         headerLogger.ServeHTTP,
		bodyLogger:   bodyLogger,
		headerLogger: headerLogger,
	}

	return logger
}

func (l *responseLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.next(w, r)
}

func (l *responseLogger) SetWriter(writer io.Writer) *responseLogger {
	l.bodyLogger.SetWriter(writer)
	l.headerLogger.SetWriter(writer)
	return l
}

type responseHeaderLogger struct {
	next   http.HandlerFunc
	writer io.Writer
}

func newResponseHeaderLogger(handlerFunc http.HandlerFunc) *responseHeaderLogger {
	logger := &responseHeaderLogger{
		next:   handlerFunc,
		writer: os.Stdout,
	}
	return logger
}

func (l *responseHeaderLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.next(w, r)
}

func (l *responseHeaderLogger) SetWriter(writer io.Writer) *responseHeaderLogger {
	l.writer = writer
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
	l.next(w, r)
}

func (l *responseBodyLogger) SetWriter(writer io.Writer) *responseBodyLogger {
	l.writer = writer
	return l
}
