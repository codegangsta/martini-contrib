package validation

import (
	"github.com/codegangsta/martini"
	"net/http"
	"reflect"
)

func tester() string {
	return "blah"
}

func Validate(obj interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		typ := reflect.TypeOf(obj).Elem()
		val := reflect.ValueOf(obj)
		errors := make(Errors)

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			validateMethod := val.MethodByName("Validate" + field.Name)
			if !validateMethod.IsValid() {
				continue
			}

			validationResult := validateMethod.Call([]reflect.Value{})[0].String()

			if validationResult != "" {
				errors[validateMethod.String()] = validationResult
			}
		}

		context.Map(errors)
	}
}

type Errors map[string]string
