package password

import (
	"crypto/rand"
	"strconv"
	"unicode/utf8"
)

type SixDigitGenerator struct {
}

func (SixDigitGenerator) GenerateCode() (string, error) {
	b := make([]byte, 2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	a1 := int(b[0]) * 1000
	a2 := int(b[1])
	code := strconv.FormatInt(int64(a1+a2), 10)

	if utf8.RuneCountInString(code) == 5 {
		code = "0" + code
	}

	return code, nil
}
