package validation

import (
	"github.com/codegangsta/martini"
	"net/http"
	"reflect"
	"strings"
)

func Validate(obj interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		typ := reflect.TypeOf(obj).Elem()
		val := reflect.ValueOf(obj)
		errors := make(Errors)

		structName := strings.Split(typ.String(), ".")[1]
		validateMethod := val.MethodByName("Validate" + structName)

		if validateMethod.IsValid() {
			validationResult := validateMethod.Call([]reflect.Value{})[0].String()
			if validationResult != "" {
				errors[validateMethod.String()] = validationResult
			}
		}

		context.Map(errors)
	}
}

type Errors map[string]string
