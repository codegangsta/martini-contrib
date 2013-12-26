# bundle
Script bundle handler for Martini.

## Usage

### NewScriptBundle(bool, bool) *ScriptBundle
`func NewScriptBundle(wrapScope bool, minify bool) *ScriptBundle`

Creates a new `ScriptBundle`.

-	setting `wrapScope` to `true` will wrap the concatenated script with an IIFE/IFFY/SIAF
-	setting `minify` will attempt to call the [Google Closure Compiler](https://developers.google.com/closure/compiler) in order to minify the resulting scripts

### (*ScriptBundle) Compile() string
`(s *ScriptBundle) Compile() string`

Concatenates, (optionally) wraps the concatenated script, (optionally) minifies the script and yields the result.

### (*ScriptBundle) AddFiles(...string)
`(s *ScriptBundle) AddFiles(files ...string)`

Adds a list of new script files to the bundle.

### (*ScriptBundle) Handler() martini.Handler
`(s *ScriptBundle) Handler() martini.Handler`

Returns a new `martini.Handler` that serves the concatenated and minified script as a single resource.

## Example

~~~ go
import (
  "github.com/codegangsta/martini"
  "github.com/codegangsta/martini-contrib/bundle"
)

func main() {
  m := martini.Classic()

  scriptBundle := bundle.NewScriptBundle(true, true)
  scriptBundle.AddFiles(
    "public/js/one.js",
    "public/js/two.js",
    "public/js/three.js",
  )

  m.Get("/js/app.js", scriptBundle.Handler())

  m.Run()
}

~~~

## Authors
* [Frank Dumont](http://github.com/fjdumont)
