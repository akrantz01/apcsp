package main

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
)

type loggingHandler struct {
	handler http.Handler
}

func (h loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set remote address to X-Forwarded-For or CF-Connecting-IP if available
	requestPort := strings.Split(r.RemoteAddr, ":")[1]
	if r.Header.Get("CF-Connecting-IP") != "" {
		r.RemoteAddr = r.Header.Get("CF-Connecting-IP") + ":" + requestPort
	} else if r.Header.Get("X-Forwarded-For") != "" {
		r.RemoteAddr = r.Header.Get("X-Forwarded-For") + ":" + requestPort
	}

	// Add logging facilities to response writer
	logger := makeLogger(w)

	// Continue on request chain
	h.handler.ServeHTTP(logger, r)

	// Log the request
	logrus.WithFields(logrus.Fields{
		"method":         r.Method,
		"status":         logger.Status(),
		"size":           logger.Size(),
		"protocol":       r.Proto,
		"remote_address": r.RemoteAddr}).Info(r.RequestURI)
}

// Copied from gorilla/handlers
func makeLogger(w http.ResponseWriter) loggingResponseWriter {
	var logger loggingResponseWriter = &responseLogger{w: w, status: http.StatusOK}
	if _, ok := w.(http.Hijacker); ok {
		logger = &hijackLogger{responseLogger{w: w, status: http.StatusOK}}
	}
	h, ok1 := logger.(http.Hijacker)
	c, ok2 := w.(http.CloseNotifier)
	if ok1 && ok2 {
		return hijackCloseNotifier{logger, h, c}
	}
	if ok2 {
		return &closeNotifyWriter{logger, c}
	}
	return logger
}

type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}

func (l *responseLogger) Flush() {
	f, ok := l.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}

func (l *responseLogger) Push(target string, opts *http.PushOptions) error {
	p, ok := l.w.(http.Pusher)
	if !ok {
		return fmt.Errorf("responseLogger does not implement http.Pusher")
	}
	return p.Push(target, opts)
}

type commonLoggingResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	Status() int
	Size() int
}

type loggingResponseWriter interface {
	commonLoggingResponseWriter
	http.Pusher
}

type closeNotifyWriter struct {
	loggingResponseWriter
	http.CloseNotifier
}

type hijackCloseNotifier struct {
	loggingResponseWriter
	http.Hijacker
	http.CloseNotifier
}

type hijackLogger struct {
	responseLogger
}

func (l *hijackLogger) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h := l.responseLogger.w.(http.Hijacker)
	conn, rw, err := h.Hijack()
	if err == nil && l.responseLogger.status == 0 {
		// Status will be StatusSwitchingProtocols (101) if there was no error and
		// WriteHeader has not been called yet
		l.responseLogger.status = http.StatusSwitchingProtocols
	}
	return conn, rw, err
}
