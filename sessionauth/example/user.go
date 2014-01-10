package main

import (
	"../../sessionauth"
)

// MyUserModel can be any struct that represents a user in my system
type MyUserModel struct {
	Name          string `form:"name"`
	Id            int64  `form:"id"`
	authenticated bool   `form:"-"`
}

// GetAnonymousUser should generate an anonymous user model
// for all sessions. This should be an unauthenticated 0 value struct.
func GenerateAnonymousUser() sessionauth.User {
	return &MyUserModel{}
}

// Login will preform any actions that are required to make a user model
// officially authenticated.
func (u *MyUserModel) Login() {
	// Update last login time
	// Add to logged-in user's list
	// etc ...
	u.authenticated = true
}

// Logout will preform any actions that are required to completely
// logout a user.
func (u *MyUserModel) Logout() {
	// Remove from logged-in user's list
	// etc ...
	u.authenticated = false
}

func (u *MyUserModel) IsAuthenticated() bool {
	return u.authenticated
}
