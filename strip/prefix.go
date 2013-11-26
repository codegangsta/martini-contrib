package strip

import (
	"github.com/codegangsta/martini"
	"net/http"
	"strings"
)

func Prefix(prefix string) martini.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if prefix == "" {
			return
		}
		if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
			r.URL.Path = p
		} else {
			http.NotFound(w, r)
		}
	}
}
