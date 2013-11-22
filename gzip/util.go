package gzip

import (
	"compress/gzip"
	"github.com/codegangsta/martini"
)

const (
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderContentEncoding = "Content-Encoding"
	HeaderContentLength   = "Content-Length"
	HeaderVary            = "Vary"
)

type gzipResponseWriter struct {
	*gzip.Writer
	w martini.ResponseWriter
}

func (grw gzipResponseWriter) WriteHeader(code int) {
	grw.w.Header().Del(HeaderContentLength)
	grw.w.WriteHeader(code)
}

func (grw gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.Writer.Write(b)
}
