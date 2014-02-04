package namespace

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"fmt"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
)

func TestNamespace(t *testing.T) {
	fmt.Println("testing")
	for _, test := range namespaceTests {
		martini := martini.Classic()
		Namespace(martini, test.namespace, func(n *MartiniNamespace) {
			handler := func() string {
				return test.method + " " + host + test.namespace + test.path
			}
			
			installHandlers(n, test, handler)
		})

		if test.method == "ANY" {
			for _, method := range httpMethods {
				testRequestMethod(t, method, martini, test)
			}
		} else {
			testRequestMethod(t, test.method, martini, test)
		}
	}
}

func TestNamespaceWithMiddleware(t *testing.T) {
	for index, test := range namespaceTests {
		martini := martini.Classic()
		Namespace(martini, test.namespace, binding.Bind(BlogPost{}), func(n *MartiniNamespace) {
			handler := func(blogPost BlogPost) {
				assertEqual(t, index, test.bindRef, blogPost)
			}
			
			installHandlers(n, test, handler)
		})

		if test.method == "ANY" {
			for _, method := range httpMethods {
				testRequestMethod(t, method, martini, test)
			}
		} else {
			testRequestMethod(t, test.method, martini, test)
		}
	}
}

func TestNotFound(t *testing.T) {
	martini := martini.Classic()
	
	Namespace(martini, "/admin", func(n *MartiniNamespace) {
		n.Get("/", func() (int, string) {
			return 200, "I'm OK"
		})
		
		n.Post("/foo", func() (int, string) {
			return 200, "I'm OK"
		})
		
		n.NotFound(func() (int, string) {
			return 418, "I'm a teapot"
		})
	})
	
	dummyTestRequestMethod(t, "GET", martini, "/admin", "/into/the/wild", 418)
	dummyTestRequestMethod(t, "GET", martini, "", "/into/the/wild", 404)
	dummyTestRequestMethod(t, "GET", martini, "/admin", "/", 200)
	dummyTestRequestMethod(t, "POST", martini, "/admin", "/foo", 200)
}

func TestForInterference(t *testing.T) {
	martini := martini.Classic()
	
	Namespace(martini, "/admin", func(n *MartiniNamespace) {
		n.Get("/foo", func() (int, string) {
			return 200, "OK"
		})
		
		n.Post("/bar", func() (int, string) {
			return 200, "OK"
		})
	})
	
	martini.Get("/foo", func() (int, string) {
		return 201, "Created"
	})
	
	dummyTestRequestMethod(t, "GET", martini, "/admin", "/foo", 200)
	dummyTestRequestMethod(t, "GET", martini, "", "/foo", 201)
	dummyTestRequestMethod(t, "POST", martini, "/admin", "/bar", 200)
	dummyTestRequestMethod(t, "POST", martini, "", "/bar", 404)
}

func dummyTestRequestMethod(t *testing.T, method string, martini *martini.ClassicMartini, ns string, path string, expectedStatus int) {
	testRequestMethod(t, method, martini, testCase{
		method,
		ns,
		path,
		`{"title": "Foo", "text": "Bar"}`,
		&BlogPost{Title: "Foo", Text: "Bar"},
		expectedStatus,
	})
}

func testRequestMethod(t *testing.T, method string, martini *martini.ClassicMartini, test testCase) {
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(method, host + test.namespace + test.path, strings.NewReader(test.payload))
	if err != nil {
		t.Error(err)
	}
	martini.ServeHTTP(recorder, req)

	if recorder.Code != test.expectStatus {
		t.Errorf("On test case %v, got status code %d but expected %d", test, recorder.Code, test.expectStatus)
	}
}

func assertEqual(t *testing.T, index int, expected *BlogPost, got BlogPost) {
	if expected.Title != got.Title {
			t.Errorf("Title: expected=%s, got=%s in test case %d\n", expected, got, index)
	}
	if expected.Text != got.Text {
			t.Errorf("Text: expected=%s, got=%s in test case %d\n", expected, got, index)
	}
}

func installHandlers(n *MartiniNamespace, test testCase, handler martini.Handler) {
	switch test.method {
	case "GET":
		n.Get(test.path, handler)
	case "PATCH":
		n.Patch(test.path, handler)
	case "POST":
		n.Post(test.path, handler)
	case "PUT":
		n.Put(test.path, handler)
	case "DELETE":
		n.Delete(test.path, handler)
	case "OPTIONS":
		n.Options(test.path, handler)
	case "HEAD":
		n.Head(test.path, handler)
	case "ANY":
		n.Any(test.path, handler)
	}
}

type (
	testCase struct {
		method      string
		namespace   string
		path        string
		payload     string
		bindRef     *BlogPost
		expectStatus int
	}
	
	BlogPost struct {
		Title       string
		Text        string
	}
)

var (
	namespaceTests = []testCase{
		testCase{
			"GET",
			"/admin",
			"/blog/create",
			`{"title": "Foo", "text": "Bar"}`,
			&BlogPost{Title: "Foo", Text: "Bar"},
			200,
		},
		testCase{
			"PATCH",
			"/admin",
			"/blog/create",
			`{"title": "Foo", "text": "Bar"}`,
			&BlogPost{Title: "Foo", Text: "Bar"},
			200,
		},
		testCase{
			"POST",
			"/admin",
			"/blog/create",
			`{"title": "Foo", "text": "Bar"}`,
			&BlogPost{Title: "Foo", Text: "Bar"},
			200,
		},
		testCase{
			"PUT",
			"/admin",
			"/blog/create",
			`{"title": "Foo", "text": "Bar"}`,
			&BlogPost{Title: "Foo", Text: "Bar"},
			200,
		},
		testCase{
			"DELETE",
			"/admin",
			"/blog/create",
			`{"title": "Foo", "text": "Bar"}`,
			&BlogPost{Title: "Foo", Text: "Bar"},
			200,
		},
		testCase{
			"OPTIONS",
			"/admin",
			"/blog/create",
			`{"title": "Foo", "text": "Bar"}`,
			&BlogPost{Title: "Foo", Text: "Bar"},
			200,
		},
		testCase{
			"HEAD",
			"/admin",
			"/blog/create",
			`{"title": "Foo", "text": "Bar"}`,
			&BlogPost{Title: "Foo", Text: "Bar"},
			200,
		},
		testCase{
			"ANY",
			"/admin",
			"/blog/create",
			`{"title": "Foo", "text": "Bar"}`,
			&BlogPost{Title: "Foo", Text: "Bar"},
			200,
		},
	}
)

var (
	httpMethods = []string{"GET", "PATCH", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}
)

const (
	host  = "http://localhost:3000"
)
