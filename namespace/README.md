# namespace

Namespaces for Martini.

[API Reference](http://godoc.org/github.com/codegangsta/martini-contrib/namespace)



## Description

Package `namespace` provides a way to handle namespaces in Martini - because everybody likes a DRY Martini.

#### Namespace

`namespace.Namespace` is a simple wrapper around martini routing to be able to install namespaces and middlewares to use across a namespace.


## Usage

This is an example to show how to use the `namespace` package:

```go
package main

import (
   "net/http"
   
   "github.com/codegangsta/martini"
   "github.com/codegangsta/martini-contrib/namespace"
 )

func main() {
	m := martini.Classic()

	// simple usage: install an admin namespace by using namespace.Namespace
	namespace.Namespace(m, "/admin", func(n *namespace.MartiniNamespace) {
		n.Get("/blog/index", func() (int, string) {
			return 200, "Hello, ", + n.Namespace + "/blog/index!"
		})
	})
	
	// advanced usage: use handlers in the whole namespace - the handlers you pass to namespace
	// will be passed down to the created routes and installed before any additional handlers.
	// Keeps your Martini nice & DRY.
	namespace.Namespace(m, "/users", binding.Bind(User{}), func(n *namespace.MartiniNamespace) {
		n.Post("/create", func() (int, string) {
			return 200, "Hello, /admin/blog/index!"
		})
		
		// Pay Attention: The Not found handler uses a globbed route on top of your namespace.
		// Therefore, always install it last.
		n.NotFound(func() (int, string) {
			return 404, "Something Special"
		})
	})

	m.Run()
}
```

## Authors
* [Beat Richartz](https://github.com/beatrichartz)
* [Jeremy Saenz](https://github.com/codegangsta)
