# config
Martini middleware that allows you to load the application's configuration from JSON files.

[API Reference](http://godoc.org/github.com/codegangsta/martini-contrib/config)

## Usage

Once you have written your configurations to a configuration file you can tell the extension to load them.

Content of configuration file is dependency in your handler function as `config.Body`, it  will be satisified by the handler like shown in example below.

`config.Body` is fully compatible with the map[interface{}]interface{} type.


For the examples, lets assume the following `config/main.json` file:

```json
{
    "greeting": "Welcome to my file configurable application"
}
```

``` go
import (
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/config"
)

func main() {
    m := martini.Classic()

    m.Use(config.File("config/main.json"))

    m.Get("/", func(conf config.Body) string{
        return fmt.Sprintf("Greeting: %s", conf["greeting"])
    })

    m.Run()
}
```

## Authors
* [Aleksandar Diklić](http://github.com/rastasheep)
