# favicon

favicon middleware for martini.

API Reference](http://godoc.org/github.com/codegangsta/martini-contrib/favicon)

## Usage

~~~~ go
import (
  "github.com/codegangsta/martini"
  "github.com/codegangsta/martini-contrib/favicon"
)

func main() {
  m := martini.Classic()
  // Returns the specified file if `r.URL.Path` is `/favicon.ico`
  m.Use(favicon.Handler("my-favicon.ico")
  m.Run()
}
~~~

Make sure to include the Gzip middleware above other middleware that alter the response body (like the render middleware).

## Authors

* [Yann Malet](http://github.com/yml/)
