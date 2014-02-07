package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/sessions"
)

const (
	keyToken    = "oauth2_token"
	keyNextPage = "next"
)

var (
	PathLogin    = "/login"
	PathLogout   = "/logout"
	PathCallback = "/oauth2callback"
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

type Tokens interface {
	AccessToken() string
	RefreshToken() string
}

type token struct {
	tk *oauth.Token
}

func (t *token) AccessToken() string {
	return t.tk.AccessToken
}

func (t *token) RefreshToken() string {
	return t.tk.RefreshToken
}

func (t *token) Expired() bool {
	return t.tk.Expired()
}

func (t *token) Expiry() time.Time {
	return t.tk.Expiry
}

func (t *token) String() string {
	return fmt.Sprintf("%v", t.tk)
}

func Google(opts *Options) martini.Handler {
	opts.AuthUrl = "https://accounts.google.com/o/oauth2/auth"
	opts.TokenUrl = "https://accounts.google.com/o/oauth2/token"
	opts.RequestUrl = "https://www.googleapis.com/oauth2/v1/userinfo"
	return OAuth2Provider(opts)
}

func OAuth2Provider(opts *Options) martini.Handler {
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

	return func(s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			switch r.URL.Path {
			case PathLogin:
				login(transport, s, w, r)
			case PathLogout:
				logout(transport, s, w, r)
			case PathCallback:
				handleOAuth2Callback(transport, s, w, r)
			}
		}
		// Inject tokens.
		c.MapTo(unmarshallToken(s), (*Tokens)(nil))
	}
}

func login(t *oauth.Transport, s sessions.Session, w http.ResponseWriter, r *http.Request) {
	next := r.URL.Query().Get(keyNextPage)
	if s.Get(keyToken) == nil {
		// user is not logged in
		http.Redirect(w, r, t.Config.AuthCodeURL(next), 302)
		return
	}
	// no need to login, redirect to the next page
	http.Redirect(w, r, next, 302)
}

func logout(t *oauth.Transport, s sessions.Session, w http.ResponseWriter, r *http.Request) {
	next := r.URL.Query().Get(keyNextPage)
	s.Delete(keyToken)
	http.Redirect(w, r, next, 302)
}

func handleOAuth2Callback(t *oauth.Transport, s sessions.Session, w http.ResponseWriter, r *http.Request) {
	next := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	tk, _ := t.Exchange(code)
	// TODO: handle error
	// Store the credentials in the session.
	val, _ := json.Marshal(tk)
	s.Set(keyToken, val)
	http.Redirect(w, r, next, 302)
}

func unmarshallToken(s sessions.Session) (t *token) {
	if s.Get(keyToken) == nil {
		return
	}
	data := s.Get(keyToken).([]byte)
	var tk oauth.Token
	json.Unmarshal(data, &tk)
	return &token{tk: &tk}
}
