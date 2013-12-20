package cache

type Item struct {
	Key        string
	Value      []byte
	Object     interface{}
	Flags      uint32
	Expiration int32
	casid      uint64
}

type Engine interface {
	Get(key string) (*Item, error)
	Set(key string, value []byte) error
}
