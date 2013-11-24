# form
POST form parser/handler for Martini.

[API Reference](http://godoc.org/github.com/codegangsta/martini-contrib/form)

## Description
form provides a convenient way to map http POST forms to struct fields. form will conveniently inject the given struct as a service for subsequent handlers.

## Usage

~~~ go
package main

import (
   "github.com/codegangsta/martini"
   "github.com/codegangsta/martini-contrib/form"
 )

type BlogPost struct{
   Title string `form:"title,required"`
   Content string `form:"content"`
}

func main() {
  m := martini.Classic()

  m.Post("/blog", form.Form(&BlogPost{}), func(blogpost *BlogPost) string {
    return blogpost.Title
  })

  m.Run()
}
~~~

## Authors
* [Jeremy Saenz](http://github.com/codegangsta)
* [Yannik Dällenbach](http://github.com/ioboi)
