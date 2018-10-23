package hkit

import (
	"io"
	"net/http"
)

// Logger can wrap an existing http.HandlerFunc and log the
// its request/response for debugging purposes
type Logger struct {
	next           http.HandlerFunc
	requestLogger  *requestLogger
	responseLogger *responseLogger
}

// NewLogger creates and returns a *Logger which implements
// ServeHTTP and is therefore an http.Handler
func NewLogger(handlerFunc http.HandlerFunc) *Logger {
	requestLogger := newRequestLogger(handlerFunc)
	responseLogger := newResponseLogger(requestLogger.ServeHTTP)
	logger := &Logger{
		next:           responseLogger.ServeHTTP,
		requestLogger:  requestLogger,
		responseLogger: responseLogger,
	}

	return logger
}

// ServeHTTP exists to make the *Logger type an http.Handler
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.next(w, r)
}

// SetWriter sets the underlying writer to be used when logging the
// request/response. It is builder-style in a sense that will return
// the same *Logger object for convenience, so this function can be
// immediately chained to, for example, like:
//
//   loggerWithCustomWriter := NewLogger(...).SetWriter()
func (l *Logger) SetWriter(writer io.Writer) *Logger {
	l.requestLogger.SetWriter(writer)
	l.responseLogger.SetWriter(writer)
	return l
}
