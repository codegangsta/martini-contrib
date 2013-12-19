package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"strings"
)

type Servers struct {
	Address string
	Port    string
}

func prepareServers(servers []Servers) []string {
	var srvs []string
	for _, value := range servers {
		srvs = append(srvs, strings.Join([]string{value.Address, value.Port}, ":"))
	}
	return srvs
}

func NewMemcacheEngine(servers ...Servers) *MemcacheEngine {
	srvs := prepareServers(servers)
	return &MemcacheEngine{
		Client: memcache.New(srvs...),
	}
}

type MemcacheEngine struct {
	Client *memcache.Client
}

func (mc *MemcacheEngine) Get(key string) (*Item, error) {
	item, err := mc.Client.Get(key)
	return &Item{
		Key:        item.Key,
		Value:      item.Value,
		Object:     item.Object,
		Flags:      item.Flags,
		Expiration: item.Expiration,
	}, err
}

func (mc *MemcacheEngine) Set(key string, value []byte) (err error) {
	return mc.Client.Set(&memcache.Item{Key: key, Value: value})
}
