package main

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	"../../sessionauth"
	"net/http"
)

func main() {
	store := sessions.NewCookieStore([]byte("secret123"))
	m := martini.Classic()

	m.Use(render.Renderer())
	m.Use(sessions.Sessions("my_session", store))
	m.Use(sessionauth.SessionUser(GenerateAnonymousUser))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Get("/login", func(r render.Render) {
		r.HTML(200, "login", nil)
	})

	m.Post("/login", binding.Bind(MyUserModel{}), func(session sessions.Session, postedUser MyUserModel, r render.Render, req *http.Request) {
		// You should verify credentials against a database or some other mechanism at this point.
		// Then you can authenticate this session.
		err := sessionauth.AuthenticateSession(session, &postedUser)
		if err != nil {
			r.JSON(500, err)
		}

		// Back to the main page
		r.Redirect("/")
	})

	m.Get("/private", sessionauth.LoginRequired, func(r render.Render, session sessions.Session, user sessionauth.User) {
		r.HTML(200, "private", user.(*MyUserModel))
	})

	m.Get("/logout", sessionauth.LoginRequired, func(session sessions.Session, user sessionauth.User, r render.Render) {
		sessionauth.Logout(session, user)
		r.Redirect("/")
	})

	m.Run()
}
