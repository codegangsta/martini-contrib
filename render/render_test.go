package render

import (
	"github.com/codegangsta/martini"
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
	m.Use(Renderer())

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
	m.Use(Renderer("fixtures"))

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

func Test_Render_Nested_HTML(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer("fixtures"))

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

func Test_Render_Error404(t *testing.T) {
	res := httptest.NewRecorder()
	r := renderer{res, nil}
	r.Error(404)
	expect(t, res.Code, 404)
}

func Test_Render_Error500(t *testing.T) {
	res := httptest.NewRecorder()
	r := renderer{res, nil}
	r.Error(500)
	expect(t, res.Code, 500)
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
