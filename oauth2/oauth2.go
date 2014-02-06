package oauth2

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/codegangsta/martini-contrib/sessions"
)

const (
	keyUserId   = "user_id"
	keyNextPage = "next_page"
)

// Represents OAuth2 backend options.
type Options struct {
	ClientId     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string

	AuthUrl    string
	TokenUrl   string
	RequestUrl string
}

func Google(opts *Options) http.HandlerFunc {
	opts.AuthUrl = "https://accounts.google.com/o/oauth2/auth"
	opts.TokenUrl = "https://accounts.google.com/o/oauth2/token"
	opts.RequestUrl = "https://www.googleapis.com/oauth2/v1/userinfo"
	return Auth(opts)
}

func Auth(opts *Options) http.HandlerFunc {
	config := &oauth.Config{
		ClientId:     opts.ClientId,
		ClientSecret: opts.ClientSecret,
		RedirectURL:  opts.RedirectURL,
		Scope:        strings.Join(opts.Scopes, " "),
		AuthURL:      opts.AuthUrl,
		TokenURL:     opts.TokenUrl,
	}

	transport := &oauth.Transport{
		Config:    config,
		Transport: http.DefaultTransport,
	}

	return func(s sessions.Session, res http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/login":
			// If not logged in, redirect to login url
			Login(transport, res, req)
		case "/logout":
			// if logged in, remove the user and redirect to next url
			Login(transport, res, req)
		case "/oauth2callback":
			// handle code and retrieve an access token, redirect to next url
			HandleOAuth2Callback(transport, res, req)
		}
	}
}

func Login(t *oauth.Transport, s sessions.Session, w http.ResponseWriter, r *http.Request) {
	next := req.URL.Query().Get(keyNextPage)
	if s.Get(keyUserID) == "" {
		// user is not logged in
		http.Redirect(w, r, t.Config.AuthCodeURL(next), 302)
		return
	}
	// no need to login, redirect to the next page
	http.Redirect(w, r, next, 0)
}

func Logout(t *oauth.Transport, s sessions.Session, w http.ResponseWriter, r *http.Request) {
	next := req.URL.Query().Get(keyNextPage)
	s.Delete(keyUserId)
	http.Redirect(w, r, next, 0)
}

func HandleOAuth2Callback(t *oauth.Transport, s sessions.Session, res http.ResponseWriter, req *http.Request) {
}
