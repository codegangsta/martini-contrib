package serialized

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type formTest struct {
	method  string
	payload string
	ok      bool
	ref     *BlogPost
}

var formTests = []formTest{
	// bad requests
	{
		"GET",
		`{blah blah blah}`,
		false,
		&BlogPost{},
	},
	{
		"POST",
		`{blah blah blah}`,
		false,
		&BlogPost{},
	},
	{
		"PUT",
		`{blah blah blah}`,
		false,
		&BlogPost{},
	},
	{
		"DELETE",
		`{blah blah blah}`,
		false,
		&BlogPost{},
	},

	// Valid requests
	{
		"GET",
		`{"content":"Test"}`,
		true,
		&BlogPost{"", "Test"},
	},
	{
		"POST",
		`{"content":"Test", "title":"TheTitle"}`,
		true,
		&BlogPost{"TheTitle", "Test"},
	},
	{
		"PUT",
		`{"content":"Test", "title":"TheTitle"}`,
		true,
		&BlogPost{"TheTitle", "Test"},
	},
	{
		"DELETE",
		`{"content":"Test", "title":"TheTitle"}`,
		true,
		&BlogPost{"TheTitle", "Test"},
	},
}

type BlogPost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func assertEqualField(t *testing.T, fieldname string, testcasenumber int, expected interface{}, got interface{}) {
	if expected != got {
		t.Errorf("%s: expected=%s, got=%s in Testcase:%i\n", fieldname, expected, got, testcasenumber)
	}
}

func verify(test formTest, t *testing.T, index int, post *BlogPost, errors Errors) {
	if !test.ok && len(errors) == 0 {
		t.Errorf("expected DeserializationError in Testcase:%i", index)
	}
	assertEqualField(t, "Title", index, test.ref.Title, post.Title)
	assertEqualField(t, "Content", index, test.ref.Content, post.Content)

}

const route = "/blogposts/create"

func Test_JsonBlogPosts(t *testing.T) {
	for index, test := range formTests {
		recorder := httptest.NewRecorder()
		handler := func(post *BlogPost, errors Errors) { verify(test, t, index, post, errors) }

		m := martini.Classic()
		switch test.method {
		case "GET":
			m.Get(route, JSON(&BlogPost{}), handler)
		case "POST":
			m.Post(route, JSON(&BlogPost{}), handler)
		case "PUT":
			m.Put(route, JSON(&BlogPost{}), handler)
		case "DELETE":
			m.Delete(route, JSON(&BlogPost{}), handler)
		}

		req, err := http.NewRequest(test.method, route, strings.NewReader(test.payload))
		if err != nil {
			t.Error(err)
		}
		m.ServeHTTP(recorder, req)
	}
}
