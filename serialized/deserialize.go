package serialized

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"io/ioutil"
	"net/http"
)

const (
	DeserializationError string = "DeserializationError"
	ReaderError          string = "ReaderError"
)

// Available errors. Use len() to check if any errors occured.
type Errors map[string]string

// Create a new JSON handler. Errors are available via bind.Errors-Service.
func JSON(jsonStruct interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		defer req.Body.Close()
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
