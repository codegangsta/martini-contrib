package oauth2

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/sessions"
)

func Test_LoginRedirect(t *testing.T) {
	recorder := httptest.NewRecorder()
	m := martini.New()
	m.Use(sessions.Sessions("my_session", sessions.NewCookieStore([]byte("secret123"))))
	m.Use(Google(&Options{
		ClientId:     "client_id",
		ClientSecret: "client_secret",
		RedirectURL:  "refresh_url",
		Scopes:       []string{"x", "y"},
	}))

	r, _ := http.NewRequest("GET", "/login", nil)
	m.ServeHTTP(recorder, r)

	location := recorder.HeaderMap["Location"][0]
	if recorder.Code != 302 {
		t.Errorf("Not being redirected to the auth page.")
	}
	if location != "https://accounts.google.com/o/oauth2/auth?access_type=&approval_prompt=&client_id=client_id&redirect_uri=refresh_url&response_type=code&scope=x+y&state=" {
		t.Errorf("Not being redirected to the right page, %v found", location)
	}
}

func Test_Logout(t *testing.T) {
	recorder := httptest.NewRecorder()
	s := sessions.NewCookieStore([]byte("secret123"))

	m := martini.Classic()
	m.Use(sessions.Sessions("my_session", s))
	m.Use(Google(&Options{
	// no need to configure
	}))

	m.Get("/", func(s sessions.Session) {
		s.Set(keyToken, "dummy token")
	})

	m.Get("/get", func(s sessions.Session) {
		if s.Get(keyToken) != nil {
			t.Errorf("User credentials are still kept in the session.")
		}
	})

	logout, _ := http.NewRequest("GET", "/logout", nil)
	index, _ := http.NewRequest("GET", "/", nil)

	m.ServeHTTP(httptest.NewRecorder(), index)
	m.ServeHTTP(recorder, logout)

	if recorder.Code != 302 {
		t.Errorf("Not being redirected to the next page.")
	}

}

func Test_InjectedTokens(t *testing.T) {
	recorder := httptest.NewRecorder()
	m := martini.Classic()
	m.Use(sessions.Sessions("my_session", sessions.NewCookieStore([]byte("secret123"))))
	m.Use(Google(&Options{
	// no need to configure
	}))
	m.Get("/", func(tokens Tokens) string {
		return "Hello world!"
	})
	r, _ := http.NewRequest("GET", "/", nil)
	m.ServeHTTP(recorder, r)
}
