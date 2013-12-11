package favicon

import (
	"crypto/md5"
	"fmt"
	"github.com/codegangsta/martini"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func ComputeEtag(reader io.Reader) ([]byte, error) {
	hash := md5.New()
	_, err := io.Copy(hash, reader)
	return hash.Sum(nil), err
}

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

	etag, err := ComputeEtag(favicon)
	if err != nil {
		log.Fatal("An error occured while trying to to compute the Etag for the favicon", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/favicon.ico" {
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
			w.Header().Set("Etag", fmt.Sprintf("%x", etag))
			http.ServeContent(w, r, fstat.Name(), fstat.ModTime(), favicon)
		}
	}
}
