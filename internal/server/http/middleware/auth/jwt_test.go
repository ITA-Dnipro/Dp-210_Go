package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateValidateToken(t *testing.T) {
	auth := NewJwtAuth(NewMockCache(), time.Minute)

	tts := []struct {
		id      string
		expires time.Duration
		wait    time.Duration
		err     bool
	}{
		{"1", time.Second, time.Second * 2, true},
		{"2", time.Second * 2, time.Second, false},
	}
	for _, tt := range tts {
		auth.Lifetime = tt.expires
		token, err := auth.CreateToken(tt.id)
		assert.Nil(t, err)
		time.Sleep(tt.wait)

		id, err := auth.ValidateToken(token)
		assert.Equal(t, tt.err, err != nil)
		if !tt.err {
			assert.EqualValues(t, tt.id, id)
		}
	}
}

func TestCreateInvalidateToken(t *testing.T) {
	auth := NewJwtAuth(NewMockCache())

	tts := []struct {
		id  string
		err bool
	}{
		{"1", true},
		{"2", true},
	}
	for _, tt := range tts {
		auth.Lifetime = time.Second * 10
		token, err := auth.CreateToken(tt.id)
		assert.Nil(t, err)

		err = auth.InvalidateToken(tt.id)
		assert.Nil(t, err)

		id, err := auth.ValidateToken(token)
		_ = id
		assert.Equal(t, tt.err, err != nil)
	}
}

type MockCache struct {
	cache map[string]string
}

func NewMockCache() *MockCache {
	return &MockCache{make(map[string]string)}
}

func (c *MockCache) Get(key string) (string, error) {
	val, ok := c.cache[key]
	if !ok {
		return "", fmt.Errorf("no such element")
	}
	return val, nil
}

func (c *MockCache) Set(key, value string) error {
	c.cache[key] = value
	return nil
}

func (c *MockCache) Del(key string) error {
	delete(c.cache, key)
	return nil
}
