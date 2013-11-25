package validation

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateBlogPosts(t *testing.T) {
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

type BlogPost struct {
	Title   string
	Content string
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

const route = "/blogposts/create"
