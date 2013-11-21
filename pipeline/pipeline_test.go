package pipeline

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPipeline(t *testing.T) {
	validQuery, _ := http.NewRequest("GET", route+"?Title=Greeting&Content=Hello+How+Are+You", nil)
	result, response := doTest(validQuery)

	if response.Code != 400 {
		t.Errorf("Expected: HTTP 400\nActual: HTTP %d", response.Code)
	}
	if result.Title != "Greeting" && result.Content != "Hello How Are You" {
		t.Errorf("Expected: Greeting-Hello How Are You\nActual: %s-%s", result.Title, result.Content)
	}

	// valid query string
	// invalid query string
	// malformed json
	// valid json
	// invalid json
}

func doTest(req *http.Request) (*BlogPost, *httptest.ResponseRecorder) {
	response := httptest.NewRecorder()
	m := martini.Classic()
	post := &BlogPost{}
	m.Get(route, Thru(post), func() {})
	m.ServeHTTP(response, req)
	return post, response
}

type BlogPost struct {
	Title   string
	Content string
}

func (this *BlogPost) ValidateBlogPost() string {
	if len(this.Title) < 4 {
		return "Title too short"
	}
	if len(this.Content) > 1024 {
		return "Content too long"
	}
	if len(this.Content) < 10 {
		return "Content too short"
	}
	return ""
}

const route = "/blogposts/create"
