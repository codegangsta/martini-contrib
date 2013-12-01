package favicon

import (
	"fmt"
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFaviconHandler(t *testing.T) {
	m := martini.Classic()
	m.Use(Handler("favicon.ico"))
	recorder := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/favicon.ico", nil)
	if err != nil {
		t.Error(err)
	}
	m.ServeHTTP(recorder, r)
	fmt.Println("recorder", recorder.Code)
	if recorder.Code != 200 {
		t.Error("An error occured while returning the favicon.ico", recorder.Code)
	}
}
