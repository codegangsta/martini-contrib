// Martini middleware/handler to enable Access Control, such as role-based access control support, Through an flag of integer kind.
package flags

import (
	"github.com/codegangsta/martini"
	"net/http"
	"reflect"
)

// BitwiseAnd use the following rules:
// 	if flag equal 0, passed
// 	if flag bitwise and base is not equal 0, passed
//  otherwise WriteHeader(statusCode)
func BitwiseAnd(base interface{}, statusCode int) martini.Handler {
	val := reflect.ValueOf(base)
	based := val.Int()
	t := val.Type()
	return func(w http.ResponseWriter, c martini.Context) {
		v := c.Get(t)
		var flag int64
		if v.IsValid() {
			flag = v.Int()
		}
		if flag == 0 || flag&based != 0 {
			return
		}
		w.WriteHeader(statusCode)
	}
}

// Less use the following rules:
// 	if flag Less than or Equal base, passed
//  otherwise WriteHeader(statusCode)
func Less(base interface{}, statusCode int) martini.Handler {
	val := reflect.ValueOf(base)
	based := val.Int()
	t := val.Type()
	return func(w http.ResponseWriter, c martini.Context) {
		v := c.Get(t)
		var flag int64
		if v.IsValid() {
			flag = v.Int()
		}
		if flag <= based {
			return
		}
		w.WriteHeader(statusCode)
	}
}

// Great use the following rules:
// 	if flag Greater than or Equal base, passed
//  otherwise WriteHeader(statusCode)
func Great(base interface{}, statusCode int) martini.Handler {
	val := reflect.ValueOf(base)
	based := val.Int()
	t := val.Type()
	return func(w http.ResponseWriter, c martini.Context) {
		v := c.Get(t)
		var flag int64
		if v.IsValid() {
			flag = v.Int()
		}
		if flag >= based {
			return
		}
		w.WriteHeader(statusCode)
	}
}

// Forbidden is Quick method the same as BitwiseAnd(base, http.StatusForbidden)
func Forbidden(base interface{}) martini.Handler {
	return BitwiseAnd(base, http.StatusForbidden)
}
