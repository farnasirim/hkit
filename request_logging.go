package hkit

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type requestLogger struct {
	next         http.HandlerFunc
	bodyLogger   *requestBodyLogger
	headerLogger *requestHeaderLogger
}

func newRequestLogger(handlerFunc http.HandlerFunc) *requestLogger {
	bodyLogger := newRequestBodyLogger(handlerFunc)
	headerLogger := newRequestHeaderLogger(bodyLogger.ServeHTTP)
	logger := &requestLogger{
		next:         headerLogger.ServeHTTP,
		headerLogger: headerLogger,
		bodyLogger:   bodyLogger,
	}

	return logger
}

func (l *requestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.next(w, r)
}

func (l *requestLogger) SetWriter(writer io.Writer) *requestLogger {
	l.bodyLogger.SetWriter(writer)
	l.headerLogger.SetWriter(writer)
	return l
}

type requestHeaderLogger struct {
	next   http.HandlerFunc
	writer io.Writer
}

func (l *requestHeaderLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer l.next(w, r)

	l.WriteKeyValue("Method", r.Method).
		WriteKeyValue("remote address", r.RemoteAddr).
		WriteLine("")

	for name, values := range r.Header {
		l.WriteKeyValue(name, strings.Join(values, ","))
	}
	l.WriteLine("")

}

func (l *requestHeaderLogger) SetWriter(writer io.Writer) *requestHeaderLogger {
	l.writer = writer
	return l
}

func (l *requestHeaderLogger) WriteKeyValue(key, value string) *requestHeaderLogger {
	return l.WriteLine(fmt.Sprintf("%s: %s", key, value))
}

func (l *requestHeaderLogger) WriteLine(s string) *requestHeaderLogger {
	l.writer.Write(append([]byte(s), []byte("\n")...))
	return l
}

func newRequestHeaderLogger(handlerFunc http.HandlerFunc) *requestHeaderLogger {
	logger := &requestHeaderLogger{
		next:   handlerFunc,
		writer: os.Stdout,
	}

	return logger
}

type requestBodyLogger struct {
	next   http.HandlerFunc
	writer io.Writer
}

func newRequestBodyLogger(handlerFunc http.HandlerFunc) *requestBodyLogger {
	logger := &requestBodyLogger{
		next:   handlerFunc,
		writer: os.Stdout,
	}

	return logger
}

func (l *requestBodyLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer l.next(w, r)

	contentBytes, duplicatedCloser, err := readAndRecreateReadCloser(r.Body)
	if err != nil {
		return
	}
	r.Body = duplicatedCloser
	l.writer.Write(contentBytes)
	l.writer.Write([]byte("\n\n"))
}

func (l *requestBodyLogger) SetWriter(writer io.Writer) *requestBodyLogger {
	l.writer = writer
	return l
}

func insertBeforeLines(toInsert []byte, container []byte) []byte {
	ret := []byte(string(toInsert))
	for _, character := range container {
		if character == '\n' {
			ret = append(ret, toInsert...)
		}
		ret = append(ret, character)
	}
	return ret
}
