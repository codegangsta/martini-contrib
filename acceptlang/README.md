# Martini Accept-Language HTTP header parsing handler

## Introduction

Using the `acceptlang` handler you can automatically parse the `Accept-Language` HTTP header and expose it as an `AcceptLanguages` struct in your handler functions. The `AcceptLanguages` struct is a slice of `AcceptLanguage` values, which contain all qualified (or unqualified) languages as were configured by the browser. The values in the slice are sorted descending by qualification (the most qualified languages will appear first).

Unqualified languages are interpreted as having a maximum qualification (1), as is defined in the HTTP/1.1 specification.

For more information:
* [HTTP/1.1 Accept-Language specification](http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.4) 

## Installation

* Install the `acceptlang` package:
    go get github.com/codegangsta/martini-contrib/acceptlang

* Import the `acceptlang` package in your code:
    import github.com/codegangsta/martini-contrib/acceptlang

## Usage

To use the handler, simply add a new handler function instance to your 
handler chain using the `acceptlang.Languages()` function as well as an 
`AcceptLanguages` dependency in your handler function. The `AcceptLanguages` dependency will be satisified by the handler.

For example:

```go
func main() {
    m := martini.Classic()

    m.Get("/", acceptlang.Languages(), func(languages acceptlang.AcceptLanguages) string {
        return fmt.Sprintf("Languages: %s", languages)
    })

    http.ListenAndServe(":8090", m)
}
```
