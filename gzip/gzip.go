package gzip

import (
	"compress/gzip"
	"github.com/codegangsta/martini"
	"net/http"
	"strings"
)

var serveGzip = func(w http.ResponseWriter, r *http.Request, c martini.Context) {
	if !strings.Contains(r.Header.Get(HeaderAcceptEncoding), "gzip") {
		return
	}
	headers := w.Header()
	headers.Set(HeaderContentEncoding, "gzip")
	headers.Set(HeaderVary, HeaderAcceptEncoding)
	gz := gzip.NewWriter(w)
	defer gz.Close()
	mw := martini.NewResponseWriter(w)
	gzw := gzipResponseWriter{gz, mw}
	c.MapTo(gzw, (*martini.ResponseWriter)(nil))
	c.Next()
}

// All returns a Handler that adds gzip compression to all requests
func All() martini.Handler {
	return serveGzip
}
