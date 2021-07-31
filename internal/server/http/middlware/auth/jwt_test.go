package middlware

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateValidateToken(t *testing.T) {
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
		token, err := CreateToken(tt.id, tt.expires)
		assert.Nil(t, err)
		time.Sleep(tt.wait)

		id, err := ValidateToken(token)
		assert.Equal(t, tt.err, err != nil)
		if !tt.err {
			assert.EqualValues(t, tt.id, id)
		}
	}
}
