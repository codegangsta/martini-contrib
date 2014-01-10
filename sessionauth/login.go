// Package login is a middleware for Martini that provides a simple way to track user sessions
// in on a website. Please see the README for a more detailed description of the package.
package login

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	"log"
	"net/http"
)

// User defines all the functions necessary to work with the user's authentication.
// The caller should implement these functions for whatever system of authentication
// they choose to use
type User interface {
	// Return whether this user is logged in or not
	IsAuthenticated() bool

	// Set any flags or extra data that should be available
	Login()

	// Clear any sensitive data out of the user
	Logout()
}

// SessionUser will try to read a valid user object out of the session. Then it will
// inject that object, or the zero value user object (from newUser) into the context.
// The newUser() function should provide a valid 0value structure for the caller's
// user type.
func SessionUser(newUser func() User) martini.Handler {
	return func(s sessions.Session, c martini.Context, l *log.Logger) {
		userJson := s.Get("AUTHUSER")
		user := newUser()

		if userJson != nil {
			err := json.Unmarshal(userJson.([]byte), user)
			if err != nil {
				l.Printf("Could not unmarshal user: %v", userJson)
			} else {
				user.Login()
			}
		}

		c.MapTo(user, (*User)(nil))
	}
}

// AuthenticatSession will mark the session and user object as authenticated. Then
// the Login() user function will be called. This function should be called after
//you have validated a user.
func AuthenticateSession(s sessions.Session, user User) error {
	user.Login()
	return UpdateUser(s, user)
}

// Logout will clear out the session and call the Logout() user function.
func Logout(s sessions.Session, user User) {
	user.Logout()
	s.Delete("AUTHUSER")
}

// LoginRequired verifies that the current user is authenticated. Any routes that
// require a login should have this handler placed in the flow. If the user is not
// authenticated, they will be redirected to /login with the "next" get parameter
// set to the attempted URL.
func LoginRequired(r render.Render, user User, w http.ResponseWriter, req *http.Request) {
	if user.IsAuthenticated() == false {
		path := fmt.Sprintf("/login?next=%s", req.URL.Path)
		r.Redirect(path, 302)
	}
}

// UpdateUser updates the User object stored in the session. This is useful incase a change
// is made to the user model that needs to persist across requests.
func UpdateUser(s sessions.Session, user User) error {
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	s.Set("AUTHUSER", userJson)
	return nil
}
