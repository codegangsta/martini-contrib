package flags

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	rolePassing = 0
	roleRobot   = 1
	roleSignOn  = 2
	roleAdmin   = 4 | roleSignOn
)

func TestFlags(t *testing.T) {
	m := martini.Classic()
	m.Use(func(r *http.Request, c martini.Context) {
		// something
		c.Map(roleRobot) // forbidden access profile
	})
	m.Get("/profile", Forbidden(roleSignOn), func() string {
		return "hello"
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/profile", nil)
	m.ServeHTTP(w, r)
	if w.Code != http.StatusForbidden {
		t.Fatal("Expected Forbidden, but got:", http.StatusText(w.Code))
	}
}
