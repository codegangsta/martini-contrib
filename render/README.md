# render
Martini middleware/handler for easily rendering serialized JSON and HTML template responses.

[API Reference](http://godoc.org/github.com/codegangsta/martini-contrib/render)

## Usage
render uses Go's [html/template](http://golang.org/pkg/html/template/) package to render html templates.

~~~ go
// main.go
package main

import (
  "github.com/codegangsta/martini"
  "github.com/codegangsta/martini-contrib/render"
)

func main() {
  m := martini.Classic()
  // render html templates from templates directory
  m.Use(render.Renderer("templates")) 

  m.Get("/", func(r render.Render) {
    r.HTML(200, "hello.tmpl", "jeremy")
  })

  m.Run()
}

~~~

~~~ html
<!-- templates/hello.tmpl -->
<h2>Hello {{.}}!</h2>
~~~

## Authors
* [Jeremy Saenz](http://github.com/codegangsta)
