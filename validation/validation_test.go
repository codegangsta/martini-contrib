package validation

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateBlogPosts(t *testing.T) {
	recorder := httptest.NewRecorder()
	handler := func(errors Errors) {
		if len(errors) == 0 {
			t.Error("Expected at least one error")
		}
	}

	m := martini.Classic()
	m.Get(route, Validate(&BlogPost{"", "..."}), handler)

	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(recorder, req)
}

type BlogPost struct {
	Title   string
	Content string
}

func (this *BlogPost) ValidateTitle() string {
	return ""
}

const route = "/blogposts/create"
