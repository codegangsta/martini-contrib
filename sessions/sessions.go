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

func Sessions(name string, store Store) martini.Handler {
	return func(res http.ResponseWriter, r *http.Request, c martini.Context, l *log.Logger) {
		session, err := store.Get(r, name)
		check(err, l)

		c.Next()

		// save session after other handlers are run
		err = session.Save(r, res)
		check(err, l)
	}
}

func check(err error, l *log.Logger) {
	if err != nil {
		l.Printf(errorFormat, err)
	}
}
