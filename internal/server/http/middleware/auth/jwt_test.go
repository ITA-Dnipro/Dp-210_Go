package auth

import (
	"testing"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"

	"github.com/stretchr/testify/assert"
)

func TestCreateValidateToken(t *testing.T) {
	auth, err := NewAuthJwt()
	if err != nil {
		t.Fatal(err)
	}

	tts := []struct {
		id      string
		role    role.Role
		expires time.Duration
		wait    time.Duration
		err     bool
	}{
		{"1", role.Viewer, time.Second, time.Second * 2, true},
		{"2", role.Operator, time.Second * 2, time.Second, false},
	}
	for _, tt := range tts {
		token, err := auth.CreateToken(UserAuth{Id: tt.id, Role: tt.role}, tt.expires)
		assert.Nil(t, err)
		time.Sleep(tt.wait)

		u, err := auth.ValidateToken(token)
		assert.Equal(t, tt.err, err != nil)
		if !tt.err {
			assert.EqualValues(t, tt.id, u.Id)
		}
	}
}
