package binding

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

/*
	To the land of Middle-ware Earth:
		One func to rule them all,
		One func to find them,
		One func to bring them all,
		And in this package BIND them.
			- Sincerely, Sauron
*/
func Bind(obj interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request, resp http.ResponseWriter) {
		var errs Errors

		contentType := req.Header.Get("Content-Type")

		if strings.Contains(contentType, "form-urlencoded") {
			context.Invoke(Form(obj))
		} else if strings.Contains(contentType, "json") {
			context.Invoke(Json(obj))
		} else {
			context.Invoke(Form(obj))
			errs = getErrors(context)
			if len(errs) > 0 {
				context.Invoke(Json(obj))
			}
		}

		context.Invoke(Validate(obj))
		errs = getErrors(context)

		if len(errs) > 0 {
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write([]byte{})
			return
		}

	}
}

func getErrors(context martini.Context) Errors {
	return context.Get(errsType).Interface().(Errors)
}

func Form(formStruct interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		req.ParseForm()
		typ := reflect.TypeOf(formStruct).Elem()
		errors := make(Errors)

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if tag := field.Tag.Get("form"); tag != "" {
				args := strings.Split(tag, ",")
				if len(args) > 0 {
					name := args[0]
					val := req.Form.Get(name)
					reflect.ValueOf(formStruct).Elem().Field(i).SetString(val)
				}
			}
		}
		context.Map(errors)
		context.Map(formStruct)
	}
}

func Json(jsonStruct interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		if req.Body != nil {
			defer req.Body.Close()
		}
		errors := make(Errors)

		content, err := ioutil.ReadAll(req.Body)
		if err != nil {
			errors["ReaderError"] = err.Error()
		} else if err = json.Unmarshal(content, jsonStruct); err != nil {
			errors[DeserializationError] = err.Error()
		}

		context.Map(errors)
		context.Map(jsonStruct)
	}
}

func Validate(obj interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		typ := reflect.TypeOf(obj).Elem()
		errors := make(Errors)

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if hasRequired(string(field.Tag)) && !reflect.ValueOf(field).IsValid() {
				errors[field.Name] = RequireError
			}
		}

		if validator, ok := obj.(Validator); ok {
			validErrors := validator.Validate()
			for key, val := range validErrors {
				if _, alreadyHasKey := errors[key]; !alreadyHasKey {
					errors[key] = val
				}
			}
		}

		context.Map(errors)
	}
}

func hasRequired(tag string) bool {
	word, required := "", "required"
	skip := false

	for i := 0; i < len(tag); i++ {
		char := tag[i]
		letter := tag[i : i+1]

		if char == '"' {
			skip = !skip
		}

		if skip {
			continue
		} else if char == ' ' || char == '\t' || char == ':' { // `required:"whatever"` will still return true
			if word == required {
				return true
			}
			word = ""
		} else {
			word += letter
		}

		if i == len(tag)-1 {
			if word == required {
				return true
			}
		}
	}

	return false
}

type (
	Errors    map[string]string
	Validator interface {
		Validate() Errors
	}
)

var (
	errsType = reflect.TypeOf(make(Errors))
)

const (
	RequireError         string = "RequireError"
	DeserializationError string = "DeserializationError"
	ReaderError          string = "ReaderError"
)
