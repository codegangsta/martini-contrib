# Contributed Martini Handlers and Utilities [![Build Status](https://drone.io/github.com/codegangsta/martini-contrib/status.png)](https://drone.io/github.com/codegangsta/martini-contrib/latest)

This package includes a variety of add-on components for Martini, a classy web framework for Go:

Install the package (**go 1.1** and greater is required):
~~~
go get github.com/codegangsta/martini-contrib
~~~

Join the [Mailing list](https://groups.google.com/forum/#!forum/martini-go)

## Available Components
* [auth](http://godoc.org/github.com/codegangsta/martini-contrib/auth) - Handlers for authentication.
* [form](http://godoc.org/github.com/codegangsta/martini-contrib/form) - Handler for parsing and mapping form fields.
* [gzip](http://godoc.org/github.com/codegangsta/martini-contrib/gzip) - Handler for adding gzip compress to requests
* [acceptlang](http://godoc.org/github.com/codegangsta/martini-contrib/acceptlang) - Handler for parsing the `Accept-Language` HTTP header.

## Examples
Want to post an example to put on the readme? Put up a Pull Request!

### Accept-Language HTTP header parsing

Using the `acceptlang` handler(s) you can automatically parse the `Accept-Language` HTTP header and expose it as an `AcceptLanguages` struct. To
use it, you add a handler to your handler chain using the `Languages()` function and define a dependency:

```go
m.Get("/", Languages(), func(languages AcceptLanguages) {
    // Use languages here when you need it
})
```

The handler implementation respects the [HTTP/1.1 Accept-Language](http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.4) specification.

## Contributing
Feel free to submit patches or file issues via GitHub. If you have an idea for a handler put up a Pull Request and we will find where it fits best!
