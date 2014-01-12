package config

import (
	"github.com/codegangsta/martini"
	"net/http"
)

type Body map[interface{}]interface{}

func File(path string) martini.Handler {
	return func(context martini.Context, request *http.Request) {
		body := make(Body)

		context.Map(body)
	}
}
