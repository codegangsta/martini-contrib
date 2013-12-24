package binding

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/codegangsta/martini"
)

func TestBind(t *testing.T) {
	index := 0
	for test, expectStatus := range bindTests {
		recorder := httptest.NewRecorder()
		handler := func(post BlogPost, errors Errors) { handle(test, t, index, post, errors) }

		m := martini.Classic()
		switch test.method {
		case "GET":
			m.Get(route, Bind(BlogPost{}), handler)
		case "POST":
			m.Post(route, Bind(BlogPost{}), handler)
		}

		req, err := http.NewRequest(test.method, test.path, strings.NewReader(test.payload))
		req.Header.Add("Content-Type", test.contentType)

		if err != nil {
			t.Error(err)
		}
		m.ServeHTTP(recorder, req)

		if recorder.Code != expectStatus {
			t.Errorf("On test case %v, got status code %d but expected %d", test, recorder.Code, expectStatus)
		}

		index++
	}
}

func TestForm(t *testing.T) {
	for index, test := range formTests {
		recorder := httptest.NewRecorder()
		handler := func(post BlogPost, errors Errors) {
			handle(test, t, index, post, errors)
		}

		m := martini.Classic()
		switch test.method {
		case "GET":
			m.Get(route, Form(BlogPost{}), handler)
		case "POST":
			m.Post(route, Form(BlogPost{}), handler)
		}

		req, err := http.NewRequest(test.method, test.path, nil)
		if err != nil {
			t.Error(err)
		}
		m.ServeHTTP(recorder, req)
	}
}

func TestJson(t *testing.T) {
	for index, test := range jsonTests {
		recorder := httptest.NewRecorder()
		handler := func(post BlogPost, errors Errors) { handle(test, t, index, post, errors) }

		m := martini.Classic()
		switch test.method {
		case "GET":
			m.Get(route, Json(BlogPost{}), handler)
		case "POST":
			m.Post(route, Json(BlogPost{}), handler)
		case "PUT":
			m.Put(route, Json(BlogPost{}), handler)
		case "DELETE":
			m.Delete(route, Json(BlogPost{}), handler)
		}

		req, err := http.NewRequest(test.method, route, strings.NewReader(test.payload))
		if err != nil {
			t.Error(err)
		}
		m.ServeHTTP(recorder, req)
	}
}

func handle(test testCase, t *testing.T, index int, post BlogPost, errors Errors) {
	assertEqualField(t, "Title", index, test.ref.Title, post.Title)
	assertEqualField(t, "Content", index, test.ref.Content, post.Content)
	assertEqualField(t, "Views", index, test.ref.Views, post.Views)

	if test.ok && errors.Count() > 0 {
		t.Errorf("%v should be OK (0 errors), but had errors: %v", test, errors)
	} else if !test.ok && errors.Count() == 0 {
		t.Errorf("%v should have errors, but was OK (0 errors): %v", test)
	}
}

func assertEqualField(t *testing.T, fieldname string, testcasenumber int, expected interface{}, got interface{}) {
	if expected != got {
		t.Errorf("%s: expected=%s, got=%s in test case %d\n", fieldname, expected, got, testcasenumber)
	}
}

func TestValidate(t *testing.T) {
	handlerMustErr := func(errors Errors) {
		if errors.Count() == 0 {
			t.Error("Expected at least one error, got 0")
		}
	}
	handlerNoErr := func(errors Errors) {
		if errors.Count() > 0 {
			t.Error("Expected no errors, got", errors.Count())
		}
	}

	performValidationTest(&BlogPost{"", "...", 0, 0}, handlerMustErr, t)
	performValidationTest(&BlogPost{"Good Title", "Good content", 0, 0}, handlerNoErr, t)

	performValidationTest(&User{Name: "Jim", Home: Address{"", ""}}, handlerMustErr, t)
	performValidationTest(&User{Name: "Jim", Home: Address{"required", ""}}, handlerNoErr, t)
}

func performValidationTest(data interface{}, handler func(Errors), t *testing.T) {
	recorder := httptest.NewRecorder()
	m := martini.Classic()
	m.Get(route, Validate(data), handler)

	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		t.Error("HTTP error:", err)
	}

	m.ServeHTTP(recorder, req)
}

func (self BlogPost) Validate(errors *Errors, req *http.Request) {
	if len(self.Title) < 4 {
		errors.Fields["Title"] = "Too short; minimum 4 characters"
	}
	if len(self.Content) > 1024 {
		errors.Fields["Content"] = "Too long; maximum 1024 characters"
	}
	if len(self.Content) < 5 {
		errors.Fields["Content"] = "Too short; minimum 5 characters"
	}
}

type (
	testCase struct {
		method      string
		path        string
		payload     string
		contentType string
		ok          bool
		ref         *BlogPost
	}

	BlogPost struct {
		Title    string `form:"title" json:"title" binding:"required"`
		Content  string `form:"content" json:"content"`
		Views    int    `form:"views" json:"views"`
		internal int    `form:"-"`
	}

	User struct {
		Name string  `json:"name" binding:"required"`
		Home Address `json:"address" binding:"required"`
	}

	Address struct {
		Street1 string `json:"street1" binding:"required"`
		Street2 string `json:"street2"`
	}
)

var (
	bindTests = map[testCase]int{
		// These should bail at the deserialization/binding phase
		testCase{
			"POST",
			path,
			`{ bad JSON `,
			"application/json",
			false,
			new(BlogPost),
		}: http.StatusBadRequest,
		testCase{
			"POST",
			path,
			`not URL-encoded but has content-type`,
			"x-www-form-urlencoded",
			false,
			new(BlogPost),
		}: http.StatusBadRequest,
		testCase{
			"POST",
			path,
			`no content-type and not URL-encoded or JSON"`,
			"",
			false,
			new(BlogPost),
		}: http.StatusBadRequest,

		// These should deserialize, then bail at the validation phase
		testCase{
			"GET",
			path + "?content=This+is+the+content",
			``,
			"x-www-form-urlencoded",
			false,
			&BlogPost{Title: "", Content: "This is the content"},
		}: http.StatusBadRequest,
		testCase{
			"GET",
			path + "",
			`{"content":"", "title":"Blog Post Title"}`,
			"application/json",
			false,
			&BlogPost{Title: "Blog Post Title", Content: ""},
		}: http.StatusBadRequest,

		// These should succeed
		testCase{
			"GET",
			path + "",
			`{"content":"This is the content", "title":"Blog Post Title"}`,
			"application/json",
			true,
			&BlogPost{Title: "Blog Post Title", Content: "This is the content"},
		}: http.StatusOK,
		testCase{
			"GET",
			path + "?content=This is the content&title=Blog+Post+Title",
			``,
			"",
			true,
			&BlogPost{Title: "Blog Post Title", Content: "This is the content"},
		}: http.StatusOK,
		testCase{
			"GET",
			path + "?content=This is the content&title=Blog+Post+Title",
			`{"content":"This is the content", "title":"Blog Post Title"}`,
			"",
			true,
			&BlogPost{Title: "Blog Post Title", Content: "This is the content"},
		}: http.StatusOK,
		testCase{
			"GET",
			path + "",
			`{"content":"This is the content", "title":"Blog Post Title"}`,
			"",
			true,
			&BlogPost{Title: "Blog Post Title", Content: "This is the content"},
		}: http.StatusOK,
	}

	formTests = []testCase{
		{
			"GET",
			path + "?content=This is the content",
			"",
			"",
			false,
			&BlogPost{Title: "", Content: "This is the content"},
		},
		{
			"POST",
			path + "?content=This is the content&title=Blog+Post+Title&views=3",
			"",
			"",
			false, // false because POST requests should have a body, not just a query string
			&BlogPost{Title: "Blog Post Title", Content: "This is the content", Views: 3},
		},
	}

	jsonTests = []testCase{
		// bad requests
		{
			"GET",
			"",
			`{blah blah blah}`,
			"",
			false,
			&BlogPost{},
		},
		{
			"POST",
			"",
			`{asdf}`,
			"",
			false,
			&BlogPost{},
		},
		{
			"PUT",
			"",
			`{blah blah blah}`,
			"",
			false,
			&BlogPost{},
		},
		{
			"DELETE",
			"",
			`{;sdf _SDf- }`,
			"",
			false,
			&BlogPost{},
		},

		// Valid-JSON requests
		{
			"GET",
			"",
			`{"content":"This is the content"}`,
			"",
			false,
			&BlogPost{Title: "", Content: "This is the content"},
		},
		{
			"POST",
			"",
			`{}`,
			"application/json",
			false,
			&BlogPost{Title: "", Content: ""},
		},
		{
			"POST",
			"",
			`{"content":"This is the content", "title":"Blog Post Title"}`,
			"",
			true,
			&BlogPost{Title: "Blog Post Title", Content: "This is the content"},
		},
		{
			"PUT",
			"",
			`{"content":"This is the content", "title":"Blog Post Title"}`,
			"",
			true,
			&BlogPost{Title: "Blog Post Title", Content: "This is the content"},
		},
		{
			"DELETE",
			"",
			`{"content":"This is the content", "title":"Blog Post Title"}`,
			"",
			true,
			&BlogPost{Title: "Blog Post Title", Content: "This is the content"},
		},
	}
)

const (
	route = "/blogposts/create"
	path  = "http://localhost:3000" + route
)
