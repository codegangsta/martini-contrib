# access-flags

[API Reference](http://gowalker.org/github.com/codegangsta/martini-contrib/access-flags)

## Description

packcage `access-flags` is Martini middleware/handler to enable Access Control, such as role-based access control support, Through an flag of integer kind.

## Usage

~~~go
package main

import (
	"github.com/Archs/martini-contrib/access-flags"
	"github.com/codegangsta/martini"
)

const (
	rolePassAll = 0
	roleRobot   = 1
	roleSignOn  = 2
	roleAdmin   = 4 | roleSignOn
)

func main() {
	m := martini.Classic()

	m.Use(judge) // first Map flag for role
	m.Get("/profile", flags.Forbidden(roleSignOn), profileHandler)
	m.Get("/admin", flags.Forbidden(roleAdmin), adminHandler)

	m.Run()
}

func judge(c martini.Context) martini.Handler{
	// something
	c.Map(roleRobot) // roleSignOn...
}

func adminHandler() string{
	// something
	return "hello"
}
func profileHandler() string{
	// something
	return "profile"
}
~~~

## Authors
* [Yu HengChun](http://github.com/achun)
