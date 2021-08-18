package auth

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MFswDQYJKoZIhvcNAQEBBQADSgAwRwJAcZOqxcmjKO7apo+jX5J8jDnkRhRRh25o
98y/iIXFdFWaAUDmr3zp60rtpgcYISske3yPhJwwqh6p8VXW3Sds1wIDAQAB
-----END PUBLIC KEY-----`

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIBOgIBAAJAcZOqxcmjKO7apo+jX5J8jDnkRhRRh25o98y/iIXFdFWaAUDmr3zp
60rtpgcYISske3yPhJwwqh6p8VXW3Sds1wIDAQABAkBDh9CXT5/iu7poJKm4Lso9
OkK/ZF9hjkV9aVFM5HUWCPCQ0mwnz00xFqntTmYYT6NbC5S5zGIA9OoTej1DhyDh
AiEAtlqIqcr7ZMiLJfD+FjBFTVPY8sfuiNKrGwoI1ai127ECIQCfclOjBDb2Aav9
/8Gh1CnzSBmoTb+3iMe93GLy9p7bBwIhAI32q4BsWwyaJ+Iw3M7PY5SQ20wfJG/2
emkBheE4h+PxAiEAnRdtsanAYKYLB0hJRSCcaDW8GaboYXIgoT2WO5yhrFcCIBkg
URG/h+mR4G6J7qPdHN2S8wK7WyqJx3TiH/nwVK+t
-----END RSA PRIVATE KEY-----`

var (
	ErrInvalidToken          = jwt.NewValidationError("invalid token", 20)
	ErrInvalidTokenStructure = jwt.NewValidationError("invalid token structure", 21)
	ErrTokenExpired          = jwt.NewValidationError("token expired", 21)
)

type JwtToken string

type Cache interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

type JwtAuth struct {
	Cache     Cache
	lifetime  time.Duration
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
}

func NewJwtAuth(cache Cache) *JwtAuth {
	auth := JwtAuth{Cache: cache}
	if err := auth.initializeAuthKeys(); err != nil {
		panic(err)
	}
	return &auth
}

func (auth *JwtAuth) initializeAuthKeys() error {
	sk, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return fmt.Errorf("initializeAuth private: %w", err)
	}
	auth.signKey = sk

	vk, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return fmt.Errorf("initializeAuth public: %w", err)
	}
	auth.verifyKey = vk

	return nil
}

func (auth *JwtAuth) CreateToken(userId string) (JwtToken, error) {
	now := time.Now()
	tomorrow := now.Add(auth.lifetime)
	claims := &AuthClaims{
		userId,
		jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: tomorrow.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t, err := token.SignedString(auth.signKey)
	return JwtToken(t), err
}

func (auth *JwtAuth) ValidateToken(t JwtToken) (string, error) {
	token, err := jwt.ParseWithClaims(string(t), &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return auth.verifyKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("validateToken: %w", err)
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return "", ErrInvalidTokenStructure
	}

	return claims.UserId, nil
}

type AuthClaims struct {
	UserId string
	jwt.StandardClaims
}
