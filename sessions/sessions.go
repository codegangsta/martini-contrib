package sessions

import (
	"github.com/codegangsta/martini"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

const (
	errorFormat = "[sessions] ERROR! %s\n"
)

type Store interface {
	sessions.Store
}

func NewCookieStore(keyPairs ...[]byte) Store {
	return sessions.NewCookieStore(keyPairs...)
}

type Session interface {
	Get(interface{}) interface{}
	Set(interface{}, interface{})
}

func Sessions(name string, store Store) martini.Handler {
	return func(res http.ResponseWriter, r *http.Request, c martini.Context, l *log.Logger) {
		s, err := store.Get(r, name)
		check(err, l)

		// Map to the Session interface
		c.MapTo(&session{s}, (*Session)(nil))

		c.Next()

		// save session after other handlers are run
		err = s.Save(r, res)
		check(err, l)
	}
}

type session struct {
	*sessions.Session
}

func (s *session) Get(key interface{}) interface{} {
	return s.Session.Values[key]
}

func (s *session) Set(key interface{}, val interface{}) {
	s.Session.Values[key] = val
}

func check(err error, l *log.Logger) {
	if err != nil {
		l.Printf(errorFormat, err)
	}
}
