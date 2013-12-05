// Package sessions contains middleware for easy session management in Martini.
//
//  package main
//
//  import (
//    "github.com/codegangsta/martini"
//    "github.com/codegangsta/martini-contrib/sessions"
//  )
//
//  func main() {
// 	  m := martini.Classic()
//
// 	  store := sessions.NewCookieStore([]byte("secret123"))
// 	  m.Use(sessions.Sessions("my_session", store))
//
// 	  m.Get("/", func(session sessions.Session) string {
// 		  session.Set("hello", "world")
// 	  })
//  }
package sessions

import (
	"github.com/codegangsta/martini"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

const (
	errorFormat = "[sessions] ERROR! %s\n"
)

// Store is an interface for custom session stores.
type Store interface {
	sessions.Store
}

// NewCookieStore returns a new CookieStore.
//
// Keys are defined in pairs to allow key rotation, but the common case is to set a single
// authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for encryption. The
// encryption key can be set to nil or omitted in the last pair, but the authentication key
// is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes. The encryption key,
// if set, must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes.
func NewCookieStore(keyPairs ...[]byte) Store {
	return sessions.NewCookieStore(keyPairs...)
}

// Session stores the values and optional configuration for a session.
type Session interface {
	// Get returns the session value associated to the given key.
	Get(key interface{}) interface{}
	// Set sets the session value associated to the given key.
	Set(key interface{}, val interface{})
	// AddFlash adds a flash message to the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	AddFlash(value interface{}, vars ...string)
	// Flashes returns a slice of flash messages from the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	Flashes(vars ...string) []interface{}
}

// Sessions is a Middleware that maps a session.Session service into the Martini handler chain.
// Sessions can use a number of storage solutions with the given store.
func Sessions(name string, store Store) martini.Handler {
	return func(res http.ResponseWriter, r *http.Request, c martini.Context, l *log.Logger) {
		// Map to the Session interface
		s := &session{name, r, l, store, nil, false}
		c.MapTo(s, (*Session)(nil))

		// clear the context, we don't need to use
		// gorilla context and we don't want memory leaks
		defer context.Clear(r)

		// Use before hook to save out the session
		rw := res.(martini.ResponseWriter)
		rw.Before(func(martini.ResponseWriter) {
			if s.Written() {
				check(s.Session().Save(r, res), l)
			}
		})
	}
}

type session struct {
	name    string
	request *http.Request
	logger  *log.Logger
	store   Store
	session *sessions.Session
	written bool
}

func (s *session) Get(key interface{}) interface{} {
	return s.Session().Values[key]
}

func (s *session) Set(key interface{}, val interface{}) {
	s.Session().Values[key] = val
	s.written = true
}

func (s *session) AddFlash(value interface{}, vars ...string) {
	s.Session().AddFlash(value, vars...)
	s.written = true
}

func (s *session) Flashes(vars ...string) []interface{} {
	s.written = true
	return s.Session().Flashes(vars...)
}

func (s *session) Session() *sessions.Session {
	if s.session == nil {
		var err error
		s.session, err = s.store.Get(s.request, s.name)
		check(err, s.logger)
	}

	return s.session
}

func (s *session) Written() bool {
	return s.written
}

func check(err error, l *log.Logger) {
	if err != nil {
		l.Printf(errorFormat, err)
	}
}
