package favicon

import (
	"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
)

// Creates a new handler that returns the favicon specified in `file` and set the `cache-control header`
func Handler(file string, maxAge int) martini.Handler {
	return func(w http.ResponseWriter, r *http.Request, log *log.Logger) {
		if r.URL.Path == "/favicon.ico" {
			log.Println("[favicon] Serving ")
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
			http.ServeFile(w, r, file)
		}
	}
}
