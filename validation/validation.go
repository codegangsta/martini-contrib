package validation

import (
	"github.com/codegangsta/martini"
	"net/http"
)

func Validate(obj Validator) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		context.Map(obj.Validate())
	}
}

type Errors map[string]string

type Validator interface {
	Validate() Errors
}
