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
	// Path to handle OAuth 2.0 logins.
	PathLogin = "/login"
	// Path to handle OAuth 2.0 logouts.
	PathLogout = "/logout"
	// Path to handle callback from OAuth 2.0 backend
	// to exchange credentials.
	PathCallback = "/oauth2callback"
	// Path to handle error cases.
	PathError = "/oauth2error"
)

// Represents OAuth2 backend options.
type Options struct {
	ClientId     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string

	AuthUrl  string
	TokenUrl string
}

// Represents a container that contains
// user's OAuth 2.0 access and refresh tokens.
type Tokens interface {
	AccessToken() string
	RefreshToken() string
	Expiry() time.Time
}

type token struct {
	tk *oauth.Token
}

// Returns the access token.
func (t *token) AccessToken() string {
	return t.tk.AccessToken
}

// Returns the refresh token.
func (t *token) RefreshToken() string {
	return t.tk.RefreshToken
}

// Returns whether the access token is
// expired or not.
func (t *token) Expired() bool {
	return t.tk.Expired()
}

// Returns the expiry time of the user's
// access token.
func (t *token) Expiry() time.Time {
	return t.tk.Expiry
}

// Formats tokens into string.
func (t *token) String() string {
	return fmt.Sprintf("%v", t.tk)
}

// Returns a new Google OAuth 2.0 backend endpoint.
func Google(opts *Options) martini.Handler {
	opts.AuthUrl = "https://accounts.google.com/o/oauth2/auth"
	opts.TokenUrl = "https://accounts.google.com/o/oauth2/token"
	return OAuth2Provider(opts)
}

// Returns a new Github OAuth 2.0 backend endpoint.
func Github(opts *Options) martini.Handler {
	opts.AuthUrl = "https://github.com/login/oauth/authorize"
	opts.TokenUrl = "https://github.com/login/oa­uth­/ac­ces­s_token"
	return OAuth2Provider(opts)
}

// Returns a generic OAuth 2.0 backend endpoint.
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
		// User is not logged in.
		fmt.Println(t.Config.AuthCodeURL(next))
		http.Redirect(w, r, t.Config.AuthCodeURL(next), 302)
		return
	}
	// No need to login, redirect to the next page.
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
	tk, err := t.Exchange(code)
	if err != nil {
		// Pass the error message, or allow dev to provide its own
		// error handler.
		http.Redirect(w, r, PathError, 302)
		return
	}
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
