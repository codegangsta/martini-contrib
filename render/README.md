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
  m.Use(render.Renderer())

  m.Get("/", func(r render.Render) {
    r.HTML(200, "hello", "jeremy")
  })

  m.Run()
}

~~~

~~~ html
<!-- templates/hello.tmpl -->
<h2>Hello {{.}}!</h2>
~~~

### Options
`render.Renderer` comes with a variety of configuration options:

~~~ go
// ...
m.Use(render.Renderer(render.Options{
  Directory: "templates", // specify what path to load the templates from
  Layout: "layout", // specify a layout template. Layouts can call {{ yield }} to render the current template.
  Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates
  Funcs: []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
  Delim: render.Delims{"{[{", "}]}"}, // Sets delimiters to the specified strings
}))
// ...
~~~

### Layouts
`render.Renderer` provides a `yield` function for layouts to access:
~~~ go
// ...
m.Use(render.Renderer(render.Options{
  Layout: "layout",
}))
// ...
~~~

~~~ html
<!-- layout.tmpl -->
<html>
  <head>
    <title>Martini Plz</title>
  </head>
  <body>
    <!-- Render the current template here -->
    {{ yield }}
  </body>
</html>
~~~

## Authors
* [Jeremy Saenz](http://github.com/codegangsta)
