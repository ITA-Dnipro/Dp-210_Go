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
	Del(key string) error
}

type JwtAuth struct {
	Cache     Cache
	Lifetime  time.Duration
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
}

func NewJwtAuth(cache Cache, lifetime time.Duration) *JwtAuth {
	auth := JwtAuth{Cache: cache, Lifetime: lifetime}
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
	tomorrow := now.Add(auth.Lifetime)
	claims := &AuthClaims{
		userId,
		jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: tomorrow.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	t, err := token.SignedString(auth.signKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	if err = auth.Cache.Set(userId, t); err != nil {
		return "", fmt.Errorf("save token for: %v; %w", userId, err)
	}

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

	if err := auth.validateInStorage(t, claims.UserId); err != nil {
		return "", err
	}

	return claims.UserId, nil
}

func (auth *JwtAuth) InvalidateToken(userId string) error {
	if err := auth.Cache.Del(userId); err != nil {
		return fmt.Errorf("invalidate token: %w", err)
	}

	return nil
}

func (auth *JwtAuth) validateInStorage(t JwtToken, userId string) error {
	tk, err := auth.Cache.Get(userId)
	if err != nil {
		return fmt.Errorf("user %v logged out", userId)
	}

	if JwtToken(tk) != t {
		return fmt.Errorf("no such token for user %v", userId)
	}

	return nil
}

type AuthClaims struct {
	UserId string
	jwt.StandardClaims
}
