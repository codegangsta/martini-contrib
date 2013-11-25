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

type Render interface {
	JSON(status int, v interface{})
	HTML(status int, name string, v interface{})
}

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
