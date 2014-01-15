package config

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Body map[interface{}]interface{}

var err error

func File(path string) martini.Handler {
	path = strings.TrimSpace(path)
	body := make(Body)

	if len(path) != 0 {
		fileContent, err := ioutil.ReadFile(path)
		if err == nil {
			err = json.Unmarshal(fileContent, &body)
		}
	}

	return func(context martini.Context, request *http.Request, log *log.Logger) {
		if err != nil {
			log.Printf("An error occurred while loading configuration file: %s", path)
		}
		context.Map(body)
	}
}
