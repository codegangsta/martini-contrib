package binding

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBind(t *testing.T) {
	index := 0
	for test, expectStatus := range bindTests {
		recorder := httptest.NewRecorder()
		handler := func(post *BlogPost, errors Errors) { handle(test, t, index, post) }

		m := martini.Classic()
		switch test.method {
		case "GET":
			m.Get(route, Bind(&BlogPost{}), handler)
		case "POST":
			m.Post(route, Bind(&BlogPost{}), handler)
		}

		req, err := http.NewRequest(test.method, route, strings.NewReader(test.payload))
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
		handler := func(post *BlogPost, errors Errors) {
			if !test.ok && len(errors) == 0 {
				t.Errorf("Expected RequireError in test case %d", index)
			}
			handle(test, t, index, post)
		}

		m := martini.Classic()
		switch test.method {
		case "GET":
			m.Get(route, Form(&BlogPost{}), handler)
		case "POST":
			m.Post(route, Form(&BlogPost{}), handler)
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
		handler := func(post *BlogPost, errors Errors) { handle(test, t, index, post) }

		m := martini.Classic()
		switch test.method {
		case "GET":
			m.Get(route, Json(&BlogPost{}), handler)
		case "POST":
			m.Post(route, Json(&BlogPost{}), handler)
		case "PUT":
			m.Put(route, Json(&BlogPost{}), handler)
		case "DELETE":
			m.Delete(route, Json(&BlogPost{}), handler)
		}

		req, err := http.NewRequest(test.method, route, strings.NewReader(test.payload))
		if err != nil {
			t.Error(err)
		}
		m.ServeHTTP(recorder, req)
	}
}

func TestTagParser(t *testing.T) {
	tests := map[string]bool{
		`form:"title" json:"title" required`: true,
		`form:"title" required json:"title"`: true,
		`required form:"title" json:"title"`: true,
		`required`:                           true,
		`form:"title" json:"title"`:          false,
		``: false,
		`form:"title" required`: true,
	}

	for input, expected := range tests {
		actual := hasRequired(input)
		if actual != expected {
			t.Errorf("Expected tag `%s` to be required=%v, but got: %v", input, expected, actual)
		}
	}
}

func handle(test testCase, t *testing.T, index int, post *BlogPost) {
	assertEqualField(t, "Title", index, test.ref.Title, post.Title)
	assertEqualField(t, "Content", index, test.ref.Content, post.Content)
}

func assertEqualField(t *testing.T, fieldname string, testcasenumber int, expected interface{}, got interface{}) {
	if expected != got {
		t.Errorf("%s: expected=%s, got=%s in test case %d\n", fieldname, expected, got, testcasenumber)
	}
}

func TestValidate(t *testing.T) {
	handlerMustErr := func(errors Errors) {
		if len(errors) == 0 {
			t.Error("Expected at least one error, got 0")
		}
	}
	handlerNoErr := func(errors Errors) {
		if len(errors) > 0 {
			t.Error("Expected no errors, got", len(errors))
		}
	}

	performValidationTest(&BlogPost{"", "..."}, handlerMustErr, t)
	performValidationTest(&BlogPost{"Good Title", "Good content"}, handlerNoErr, t)
}

func performValidationTest(post *BlogPost, handler func(Errors), t *testing.T) {
	recorder := httptest.NewRecorder()
	m := martini.Classic()
	m.Get(route, Validate(post), handler)

	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		t.Error("HTTP error:", err)
	}

	m.ServeHTTP(recorder, req)
}

func (this *BlogPost) Validate() Errors {
	errs := make(Errors)

	if len(this.Title) < 10 {
		errs["Title"] = "Title too short"
	}
	if len(this.Content) > 1024 {
		errs["Content"] = "Content too long"
	}
	if len(this.Content) < 10 {
		errs["Content"] = "Content too short"
	}

	return errs
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
		Title   string `form:"title" json:"title" required`
		Content string `form:"content" json:"content"`
	}
)

var (
	bindTests = map[testCase]int{
		// These should bail at the deserialization/binding phase
		testCase{
			"POST",
			"",
			`{ bad JSON `,
			"application/json",
			false,
			new(BlogPost),
		}: http.StatusBadRequest,
		testCase{
			"POST",
			"",
			`not URL-encoded: "see?"`,
			"x-www-form-urlencoded",
			false,
			new(BlogPost),
		}: http.StatusBadRequest,

		// These should deserialize, then bail at the validation phase
		testCase{
			"GET",
			"http://localhost:3000/blogposts/create?content=Test",
			``,
			"x-www-form-urlencoded",
			false,
			&BlogPost{"", "Test"},
		}: http.StatusBadRequest,
	}

	formTests = []testCase{
		{
			"GET",
			"http://localhost:3000/blogposts/create?content=Test",
			"",
			"",
			true,
			&BlogPost{"", "Test"},
		},
		{
			"POST",
			"http://localhost:3000/blogposts/create?content=Test&title=TheTitle",
			"",
			"",
			true,
			&BlogPost{"TheTitle", "Test"},
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
			`{blah blah blah}`,
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
			`{blah blah blah}`,
			"",
			false,
			&BlogPost{},
		},

		// Valid requests
		{
			"GET",
			"",
			`{"content":"Test"}`,
			"",
			true,
			&BlogPost{"", "Test"},
		},
		{
			"POST",
			"",
			`{"content":"Test", "title":"TheTitle"}`,
			"",
			true,
			&BlogPost{"TheTitle", "Test"},
		},
		{
			"PUT",
			"",
			`{"content":"Test", "title":"TheTitle"}`,
			"",
			true,
			&BlogPost{"TheTitle", "Test"},
		},
		{
			"DELETE",
			"",
			`{"content":"Test", "title":"TheTitle"}`,
			"",
			true,
			&BlogPost{"TheTitle", "Test"},
		},
	}
)

const route = "/blogposts/create"
