package gzip

import (
	"compress/gzip"
	"net/http"
)

const (
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderContentEncoding = "Content-Encoding"
	HeaderContentLength   = "Content-Length"
	HeaderVary            = "Vary"
)

type gzipResponseWriter struct {
	*gzip.Writer
	w http.ResponseWriter
}

func (grw gzipResponseWriter) Header() http.Header {
	return grw.w.Header()
}

func (grw gzipResponseWriter) WriteHeader(code int) {
	grw.Header().Del(HeaderContentLength)
	grw.w.WriteHeader(code)
}

func (grw gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.Writer.Write(b)
}
