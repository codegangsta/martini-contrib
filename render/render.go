// Package render is a middleware for Martini that provides easy JSON serialization and HTML template rendering.
//
//  package main
//
//  import (
//    "github.com/codegangsta/martini"
//    "github.com/codegangsta/martini-contrib/render"
//  )
//
//  func main() {
//    m := martini.Classic()
//    m.Use(render.Renderer("templates"))
//
//    m.Get("/html", func(r render.Render) {
//      r.HTML(200, "mytemplate.tmpl", nil)
//    })
//
//    m.Get("/json", func(r render.Render) {
//      r.JSON(200, "hello world")
//    })
//
//    m.Run()
//  }
package render

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"html/template"
	"net/http"
	"path/filepath"
)

const (
	ContentType = "Content-Type"
	ContentJSON = "application/json"
	ContentHTML = "text/html"
)

// Render is a service that can be injected into a Martini handler. Render provides functions for easily writing JSON and
// HTML templates out to a http Response.
type Render interface {
	// JSON writes the given status and JSON serialized version of the given value to the http.ResponseWriter.
	JSON(status int, v interface{})
	// HTML renders a html template specified by the name and writes the result and given status to the http.ResponseWriter.
	HTML(status int, name string, v interface{})
}

// Renderer is a Middleware that maps a render.Render service into the Martini handler chain. Renderer will compile templates
// globbed in the given dir. Templates must have the .tmpl extension to be compiled.
//
// If MARTINI_ENV is set to "" or "development" then templates will be recompiled on every request. For more performance, set the
// MARTINI_ENV environment variable to "production"
func Renderer(dir string) martini.Handler {
	t := compile(dir)
	return func(res http.ResponseWriter, c martini.Context) {
		// recompile for easy development
		if martini.Env == martini.Dev {
			t = compile(dir)
		}
		c.MapTo(&renderer{res, t}, (*Render)(nil))
	}
}

func compile(dir string) *template.Template {
	t, err := template.ParseGlob(filepath.Join(dir, "*.tmpl"))
	if err != nil {
		// do nothing for now?
		t = template.New("null")
	}
	return t
}

type renderer struct {
	http.ResponseWriter
	t *template.Template
}

func (r *renderer) JSON(status int, v interface{}) {
	r.Header().Set(ContentType, ContentJSON)
	r.WriteHeader(status)

	result, err := json.Marshal(v)
	if err != nil {
		http.Error(r, err.Error(), 500)
	}

	r.Write(result)
}

func (r *renderer) HTML(status int, name string, binding interface{}) {
	r.Header().Set(ContentType, ContentHTML)
	r.WriteHeader(status)
	err := r.t.ExecuteTemplate(r, name, binding)
	if err != nil {
		http.Error(r, err.Error(), 500)
	}
}
