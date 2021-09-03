package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestCreateValidateToken(t *testing.T) {
	tts := []struct {
		id      string
		role    entity.Role
		expires time.Duration
		wait    time.Duration
		err     bool
	}{
		{"1", entity.Viewer, time.Second, time.Second * 2, true},
		{"2", entity.Operator, time.Second * 2, time.Second, false},
	}
	for _, tt := range tts {
		auth, err := NewJwtAuth(NewMockCache(), tt.expires)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.CreateToken(UserAuth{Id: tt.id, Role: tt.role})
		assert.Nil(t, err)
		time.Sleep(tt.wait)

		u, err := auth.ValidateToken(token)
		assert.Equal(t, tt.err, err != nil)
		if !tt.err {
			assert.EqualValues(t, tt.id, u.Id)
		}
	}
}

func TestCreateInvalidateToken(t *testing.T) {
	auth, err := NewJwtAuth(NewMockCache(), time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	tts := []struct {
		id   string
		err  bool
		role entity.Role
	}{
		{"1", true, entity.Viewer},
		{"2", true, entity.Viewer},
	}
	for _, tt := range tts {
		auth.Lifetime = time.Second * 10
		token, err := auth.CreateToken(UserAuth{Id: tt.id, Role: tt.role})
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

func (c *MockCache) Get(ctx context.Context, key string) (string, error) {
	val, ok := c.cache[key]
	if !ok {
		return "", fmt.Errorf("no such element")
	}
	return val, nil
}

func (c *MockCache) Set(ctx context.Context, key, value string) error {
	c.cache[key] = value
	return nil
}

func (c *MockCache) Del(ctx context.Context, key string) error {
	delete(c.cache, key)
	return nil
}
