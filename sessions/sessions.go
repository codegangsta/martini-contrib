package sessions

import (
	"github.com/codegangsta/martini"
	"github.com/gorilla/sessions"
)

type Store interface {
	sessions.Store
}

func Sessions(name string, store Store) martini.Handler {
	return nil
}
