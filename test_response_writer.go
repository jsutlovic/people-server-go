package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type testResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// Don't need this yet because we get it for free:
func (w *testResponseWriter) Write(data []byte) (n int, err error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	size, err := w.ResponseWriter.Write(data)
	w.size += size
	return size, err
}

func (w *testResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *testResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w *testResponseWriter) Written() bool {
	return w.statusCode != 0
}

func (w *testResponseWriter) Size() int {
	return w.size
}

func (w *testResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (w *testResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w *testResponseWriter) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
