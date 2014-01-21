# recovery
Martini middleware that renders a pretty recovery page when your app panics.
Should replace the default martini.Recovery.

## Usage

~~~ go
package main

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/recovery"
)

func main() {
	m := martini.Classic()
	if martini.Env != martini.Prod {
		m.Use(recovery.Recovery())
	}
	m.Run()
}
~~~

## Authors
* [Andrew Wayne](http://github.com/dre1080)
