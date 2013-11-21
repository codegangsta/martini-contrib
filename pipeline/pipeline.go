package pipeline

import (
	"github.com/codegangsta/martini"
	"net/http"
	"reflect"
)

func Thru(obj interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		handler := form.Form(obj)
		handler(context, req)
		context.Get(reflect.TypeOf())
	}
}
