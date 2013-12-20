package cache

import (
	"github.com/codegangsta/martini"
	"log"
	"net/http"
)

const (
	errorFormat = "[cache] ERROR! %s\n"
)

type Cache interface {
	Get(key string) (*Item, error)
	Set(key string, value []byte) error
}

func Caches(engine Engine) martini.Handler {
	return func(res http.ResponseWriter, c martini.Context) {
		c.MapTo(&cache{engine}, (*Cache)(nil))
	}
}

type cache struct {
	engine Engine
}

func (c *cache) Get(key string) (*Item, error) {
	return c.engine.Get(key)
}

func (c *cache) Set(key string, value []byte) (err error) {
	err = c.engine.Set(key, value)
	return err
}

func check(err error) {
	if err != nil {
		log.Printf(errorFormat, err)
	}
}
