package cache

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_MemCache(t *testing.T) {
	println()
	println("==== Memcache Test ====")
	m := martini.Classic()

	engine := NewMemcacheEngine(Servers{
		"127.0.0.1", "11211",
	})

	m.Use(Caches(engine))

	m.Get("/setmemcache", func(cache Cache) string {
		cache.Set("hello", []byte("world"))
		return "OK"
	})

	m.Get("/getmemcache", func(cache Cache) string {
		if string(cache.Get("hello")) != "world" {
			t.Error("Cache writing failed")
		}
		return "OK"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/setmemcache", nil)
	m.ServeHTTP(res, req)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/getmemcache", nil)
	m.ServeHTTP(res2, req2)
}

func Test_FileCache(t *testing.T) {
	println()
	println("==== FileCache Test ====")

	m := martini.Classic()

	engine := NewFilecacheEngine("__cache__")

	m.Use(Caches(engine))

	m.Get("/setfilecache", func(cache Cache) string {
		cache.Set("hello", []byte("world"))
		return "OK"
	})

	m.Get("/getfilecache", func(cache Cache) string {
		if string(cache.Get("hello")) != "world" {
			t.Error("Cache writing failed")
		}
		return "OK"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/setfilecache", nil)
	m.ServeHTTP(res, req)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/getfilecache", nil)
	m.ServeHTTP(res2, req2)
}