# Contributed Martini Handlers and Utilities [![wercker status](https://app.wercker.com/status/6e73d91b3a2bdb85a74cd61d380248d7 "wercker status")](https://app.wercker.com/project/bykey/6e73d91b3a2bdb85a74cd61d380248d7)

This package includes a variety of add-on components for [Martini][], a classy web framework for Go:

Install all the packages (**go 1.1** and greater is required):
~~~
go get github.com/codegangsta/martini-contrib/...
~~~

Join the [Mailing list](https://groups.google.com/forum/#!forum/martini-go)

## Available Components
* [acceptlang](https://github.com/codegangsta/martini-contrib/tree/master/acceptlang) - Handler for parsing the `Accept-Language` HTTP header.
* [auth](https://github.com/codegangsta/martini-contrib/tree/master/auth) - Handlers for authentication.
* [binding](https://github.com/codegangsta/martini-contrib/tree/master/binding) - Handler for mapping/validating a raw request into a structure.
* [encoder](https://github.com/codegangsta/martini-contrib/tree/master/encoder) - Encoder service for rendering data in several formats and content negotiation.
* [gzip](https://github.com/codegangsta/martini-contrib/tree/master/gzip) - Handler for adding gzip compress to requests
* [method](https://github.com/codegangsta/martini-contrib/tree/master/method) - HTTP method overriding via Header or form fields.
* [render](https://github.com/codegangsta/martini-contrib/tree/master/render) - Handler that provides a service for easily rendering JSON and HTML templates.
* [secure](https://github.com/codegangsta/martini-contrib/tree/master/secure) - Implements a few quick security wins.
* [sessions](https://github.com/codegangsta/martini-contrib/tree/master/sessions) - Handler that provides a Session service.
* [strip](https://github.com/codegangsta/martini-contrib/tree/master/strip) - URL Prefix stripping.
* [web](https://github.com/codegangsta/martini-contrib/tree/master/web) - web.go Context compatibility.

## Contributing
Feel free to submit patches or file issues via GitHub. If you have an idea for a handler put up a Pull Request and we will find where it fits best!

### Be a Collaborator
If you want to be a maintainer of martini-contrib and get full repo access contact [Jeremy Saenz](http://github.com/codegangsta). I will automatically add you as a collaborator if you contribute a package yourself so you can fix bugs without a pull request.

[Martini]: http://martini.codegangsta.io/
