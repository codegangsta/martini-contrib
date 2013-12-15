// Package binding transforms, with validation, a raw request into
// a populated structure used by your application logic.
package binding

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
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

// Bind accepts a copy of an empty struct and populates it with
// values from the request (if deserialization is successful). It
// wraps up the functionality of the Form and Json middleware
// according to the Content-Type of the request, and it guesses
// if no Content-Type is specified. Bind invokes the ErrorHandler
// middleware to bail out if errors occurred. If you want to perform
// your own error handling, use Form or Json middleware directly.
func Bind(obj interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		contentType := req.Header.Get("Content-Type")

		if strings.Contains(contentType, "form-urlencoded") {
			context.Invoke(Form(obj))
		} else if strings.Contains(contentType, "json") {
			context.Invoke(Json(obj))
		} else {
			context.Invoke(Json(obj))
			if getErrors(context).Count() > 0 {
				context.Invoke(Form(obj))
			}
		}

		context.Invoke(ErrorHandler)
	}
}

// Form is middleware to deserialize Form-encoded data from the request.
// It gets data from the form-urlencoded payload, if present, or from the
// query string as well. It uses the http.Request.ParseForm() method to
// perform deserialization, then reflection is used to map each field
// into the struct with the proper type.
func Form(formStruct interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		ensureNotPointer(formStruct)
		formStruct := reflect.New(reflect.TypeOf(formStruct))
		errors := newErrors()
		parseErr := req.ParseForm()

		if parseErr != nil {
			errors.Overall[DeserializationError] = parseErr.Error()
		}

		typ := formStruct.Elem().Type()

		for i := 0; i < typ.NumField(); i++ {
			typeField := typ.Field(i)
			if inputFieldName := typeField.Tag.Get("form"); inputFieldName != "" {
				inputValue := req.Form.Get(inputFieldName)
				structField := formStruct.Elem().Field(i)

				if !structField.CanSet() {
					continue
				}

				setWithProperType(typeField, inputValue, structField, inputFieldName, errors)
			}
		}

		validateAndMap(formStruct, context, errors)
	}
}

// Json is middleware to deserialize a JSON payload from the request
// into the struct that is passed in. The resulting struct is then
// validated, but no error handling is actually performed here.
func Json(jsonStruct interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		ensureNotPointer(jsonStruct)
		jsonStruct := reflect.New(reflect.TypeOf(jsonStruct))
		errors := newErrors()

		if req.Body != nil {
			defer req.Body.Close()
		}

		content, err := ioutil.ReadAll(req.Body)
		if err != nil {
			errors.Overall[ReaderError] = err.Error()
		} else if err = json.Unmarshal(content, jsonStruct.Interface()); err != nil {
			errors.Overall[DeserializationError] = err.Error()
		}

		validateAndMap(jsonStruct, context, errors)
	}
}

// Validate is middleware to enforce required fields. If the struct
// passed in is a Validator, then the user-defined Validate method
// is executed, and its errors are mapped to the context. This middleware
// performs no error handling: it merely detects them and maps them.
func Validate(obj interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		errors := newErrors()
		validateStruct(errors, obj)

		if validator, ok := obj.(Validator); ok {
			validator.Validate(errors, req)
		}
		context.Map(*errors)

	}
}

func validateStruct(errors *Errors, obj interface{}) {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Allow ignored fields in the struct
		if field.Tag.Get("form") == "-" {
			continue
		}

		fieldValue := val.Field(i).Interface()
		zero := reflect.Zero(field.Type).Interface()

		if strings.Index(field.Tag.Get("binding"), "required") > -1 {
			if field.Type.Kind() == reflect.Struct {
				validateStruct(errors, fieldValue)
			} else if reflect.DeepEqual(zero, fieldValue) {
				errors.Fields[field.Name] = RequireError
			}
		}
	}
}

// ErrorHandler simply counts the number of errors in the
// context and, if more than 0, writes a 400 Bad Request
// response and a JSON payload describing the errors with
// the "Content-Type" set to "application/json".
// Middleware remaining on the stack will not even see the request
// if, by this point, there are any errors.
// This is a "default" handler, of sorts, and you are
// welcome to use your own instead. The Bind middleware
// invokes this automatically for convenience.
func ErrorHandler(errs Errors, resp http.ResponseWriter) {
	if errs.Count() > 0 {
		resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		resp.WriteHeader(http.StatusBadRequest)
		errOutput, _ := json.Marshal(errs)
		resp.Write(errOutput)
		return
	}
}

// This sets the value in a struct of an indeterminate type to the
// matching value from the request (via Form middleware) in the
// same type, so that not all deserialized values have to be strings.
// Supported types are string, int, float, and bool.
func setWithProperType(typeField reflect.StructField, val string, structField reflect.Value, nameInTag string, errors *Errors) {
	switch typeField.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val == "" {
			val = "0"
		}
		intVal, err := strconv.Atoi(val)
		if err != nil {
			errors.Fields[nameInTag] = IntegerTypeError
		} else {
			structField.SetInt(int64(intVal))
		}
	case reflect.Bool:
		if val == "" {
			val = "false"
		}
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			errors.Fields[nameInTag] = BooleanTypeError
		} else {
			structField.SetBool(boolVal)
		}
	case reflect.Float32:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			errors.Fields[nameInTag] = FloatTypeError
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.Float64:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			errors.Fields[nameInTag] = FloatTypeError
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.String:
		structField.SetString(val)
	}
}

// Don't pass in pointers to bind to. Can lead to bugs. See:
// https://github.com/codegangsta/martini-contrib/issues/40
// https://github.com/codegangsta/martini-contrib/pull/34#issuecomment-29683659
func ensureNotPointer(obj interface{}) {
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		panic("Pointers are not accepted as binding models")
	}
}

// Performs validation and combines errors from validation
// with errors from deserialization, then maps both the
// resulting struct and the errors to the context.
func validateAndMap(obj reflect.Value, context martini.Context, errors *Errors) {
	context.Invoke(Validate(obj.Interface()))
	errors.combine(getErrors(context))
	context.Map(*errors)
	context.Map(obj.Elem().Interface())
}

func newErrors() *Errors {
	return &Errors{make(map[string]string), make(map[string]string)}
}

func getErrors(context martini.Context) Errors {
	return context.Get(reflect.TypeOf(Errors{})).Interface().(Errors)
}

func (this *Errors) combine(other Errors) {
	for key, val := range other.Fields {
		if _, exists := this.Fields[key]; !exists {
			this.Fields[key] = val
		}
	}
	for key, val := range other.Overall {
		if _, exists := this.Overall[key]; !exists {
			this.Overall[key] = val
		}
	}
}

// Total errors is the sum of errors with the request overall
// and errors on individual fields.
func (self Errors) Count() int {
	return len(self.Overall) + len(self.Fields)
}

type (
	// Errors represents the contract of the response body when the
	// binding step fails before getting to the application.
	Errors struct {
		Overall map[string]string `json:"overall"`
		Fields  map[string]string `json:"fields"`
	}

	// Implement the Validator interface to define your own input
	// validation before the request even gets to your application.
	// The Validate method will be executed during the validation phase.
	Validator interface {
		Validate(*Errors, *http.Request)
	}
)

const (
	RequireError         string = "Required"
	DeserializationError string = "DeserializationError"
	ReaderError          string = "ReaderError"
	IntegerTypeError     string = "IntegerTypeError"
	BooleanTypeError     string = "BooleanTypeError"
	FloatTypeError       string = "FloatTypeError"
)
