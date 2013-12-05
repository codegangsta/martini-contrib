package sessions

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Sessions(t *testing.T) {
	m := martini.Classic()

	store := NewCookieStore([]byte("secret123"))
	m.Use(Sessions("my_session", store))

	m.Get("/testsession", func(session Session) string {
		session.Set("hello", "world")
		return "OK"
	})

	m.Get("/show", func(session Session) string {
		if session.Get("hello") != "world" {
			t.Error("Session writing failed")
		}
		return "OK"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/testsession", nil)
	m.ServeHTTP(res, req)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/show", nil)
	req2.Header.Set("Cookie", res.Header().Get("Set-Cookie"))
	m.ServeHTTP(res2, req2)
}

func Test_Flashes(t *testing.T) {
	m := martini.Classic()

	store := NewCookieStore([]byte("secret123"))
	m.Use(Sessions("my_session", store))

	m.Get("/set", func(session Session) string {
		session.AddFlash("hello world")
		return "OK"
	})

	m.Get("/show", func(session Session) string {
		l := len(session.Flashes())
		if l != 1 {
			t.Error("Flashes count does not equal 1. Equals ", l)
		}
		return "OK"
	})

	m.Get("/showagain", func(session Session) string {
		l := len(session.Flashes())
		if l != 0 {
			t.Error("flashes count is not 0 after reading. Equals ", l)
		}
		return "OK"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/set", nil)
	m.ServeHTTP(res, req)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/show", nil)
	req2.Header.Set("Cookie", res.Header().Get("Set-Cookie"))
	m.ServeHTTP(res2, req2)

	res3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/showagain", nil)
	req3.Header.Set("Cookie", res2.Header().Get("Set-Cookie"))
	m.ServeHTTP(res3, req3)
}
