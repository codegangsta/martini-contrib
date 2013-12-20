package render

import (
	"github.com/codegangsta/martini"
	"html/template"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type Greeting struct {
	One string `json:"one"`
	Two string `json:"two"`
}

func Test_Render_JSON(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
	// nothing here to configure
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.JSON(300, Greeting{"hello", "world"})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 300)
	expect(t, res.Header().Get(ContentType), ContentJSON)
	expect(t, res.Body.String(), `{"one":"hello","two":"world"}`)
}

func Test_Render_HTML(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		Directory: "fixtures/basic",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "hello", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML)
	expect(t, res.Body.String(), "<h1>Hello jeremy</h1>\n")
}

func Test_Render_Extensions(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		Directory:  "fixtures/basic",
		Extensions: []string{".tmpl", ".html"},
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "hypertext", nil)
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML)
	expect(t, res.Body.String(), "Hypertext!\n")
}

func Test_Render_Funcs(t *testing.T) {

	m := martini.Classic()
	m.Use(Renderer(Options{
		Directory: "fixtures/custom_funcs",
		Funcs: []template.FuncMap{
			{
				"myCustomFunc": func() string {
					return "My custom function"
				},
			},
		},
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "index", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Body.String(), "My custom function\n")
}

func Test_Render_Layout(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		Directory: "fixtures/basic",
		Layout:    "layout",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "content", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Body.String(), "head\n<h1>jeremy</h1>\n\nfoot\n")
}

func Test_Render_Nested_HTML(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		Directory: "fixtures/basic",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "admin/index", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML)
	expect(t, res.Body.String(), "<h1>Admin jeremy</h1>\n")
}

func Test_Render_Delimiters(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		Delims:    Delims{"{[{", "}]}"},
		Directory: "fixtures/basic",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "delims", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML)
	expect(t, res.Body.String(), "<h1>Hello jeremy</h1>")
}

func Test_Render_Error404(t *testing.T) {
	res := httptest.NewRecorder()
	r := renderer{res, nil, Options{}, ""}
	r.Error(404)
	expect(t, res.Code, 404)
}

func Test_Render_Error500(t *testing.T) {
	res := httptest.NewRecorder()
	r := renderer{res, nil, Options{}, ""}
	r.Error(500)
	expect(t, res.Code, 500)
}

func Test_Render_Default_Charset_JSON(t *testing.T) {
	m := martini.Classic()
	charset := "UTF-8"
	m.Use(Renderer(Options{
		Charset: charset,
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.JSON(300, Greeting{"hello", "world"})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 300)
	expect(t, res.Header().Get(ContentType), ContentJSON+"; charset="+charset)
	expect(t, res.Body.String(), `{"one":"hello","two":"world"}`)
}

func Test_Render_Default_Charset_HTML(t *testing.T) {
	m := martini.Classic()
	charset := "ISO-8859-1"
	m.Use(Renderer(Options{
		Directory: "fixtures/basic",
		Charset:   charset,
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "hello", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML+"; charset="+charset)
	expect(t, res.Body.String(), "<h1>Hello jeremy</h1>\n")
}

func Test_Render_Blank_Charset(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		Charset: "",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.JSON(300, Greeting{"hello", "world"})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 300)
	expect(t, res.Header().Get(ContentType), ContentJSON)
	expect(t, res.Body.String(), `{"one":"hello","two":"world"}`)
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
