package favicon

import (
	"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"path/filepath"
)

// Creates a new handler that returns the favicon specified in `file` and set the `cache-control header`
func Handler(name string, maxAge int) martini.Handler {
	directory, file := filepath.Split(name)
	dir := http.Dir(directory)
	favicon, err := dir.Open(file)
	if err != nil {
		// TODO (yml): Should I swallow this error
		log.Fatal("An error occured while trying to open the favicon", err)
	}

	fstat, err := favicon.Stat()
	if err != nil {
		// TODO (yml): Should I swallow this error
		log.Fatal("An error occured while accessing the file stats for the favicon", err)
	}

	return func(w http.ResponseWriter, r *http.Request, log *log.Logger) {
		if r.URL.Path == "/favicon.ico" {
			log.Println("[favicon] Serving ")
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
			http.ServeContent(w, r, fstat.Name(), fstat.ModTime(), favicon)
		}
	}
}
