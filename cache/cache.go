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
	Get(key string) []byte
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

func (c *cache) Get(key string) (value []byte) {
	cache_item, err := c.engine.Get(key)
	check(err)
	return cache_item.Value
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
