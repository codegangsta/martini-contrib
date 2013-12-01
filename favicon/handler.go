package favicon

import (
	"github.com/codegangsta/martini"
	"log"
	"net/http"
)

// Creates a new handler that returns the favicon specified in `file`
func Handler(file string) martini.Handler {
	return func(w http.ResponseWriter, r *http.Request, log *log.Logger) {
		if r.URL.Path == "/favicon.ico" {
			log.Println("[favicon] Serving ")
			http.ServeFile(w, r, file)
		}
	}
}
