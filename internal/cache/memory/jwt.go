package memory

import "fmt"

type JwtCache struct {
	cache map[string]string
}

func NewJwtCache() *JwtCache {
	return &JwtCache{make(map[string]string)}
}

func (c *JwtCache) Get(key string) (string, error) {
	val, ok := c.cache[key]
	if !ok {
		return "", fmt.Errorf("no such element")
	}
	return val, nil
}

func (c *JwtCache) Set(key, value string) error {
	c.cache[key] = value
	return nil
}

func (c *JwtCache) Del(key string) error {
	delete(c.cache, key)
	return nil
}
