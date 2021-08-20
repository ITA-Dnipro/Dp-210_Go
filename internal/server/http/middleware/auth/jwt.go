package auth

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"

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

type JwtToken string

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

var (
	ErrInvalidToken          = jwt.NewValidationError("invalid token", 20)
	ErrInvalidTokenStructure = jwt.NewValidationError("invalid token structure", 21)
	ErrTokenExpired          = jwt.NewValidationError("token expired", 21)
)

func init() {
	InitializeAuthKeys()
}

func InitializeAuthKeys() error {
	sk, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return fmt.Errorf("initializeAuth private: %w", err)
	}
	signKey = sk

	vk, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return fmt.Errorf("initializeAuth public: %w", err)
	}
	verifyKey = vk

	return nil
}

func CreateToken(user UserAuth, lifetime time.Duration) (JwtToken, error) {
	now := time.Now()
	tomorrow := now.Add(lifetime)
	claims := &AuthClaims{
		user.Id,
		user.Role,
		jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: tomorrow.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t, err := token.SignedString(signKey)
	return JwtToken(t), err
}

func ValidateToken(t JwtToken) (UserAuth, error) {
	token, err := jwt.ParseWithClaims(string(t), &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil {
		return UserAuth{}, fmt.Errorf("validateToken: %w", err)
	}

	if !token.Valid {
		return UserAuth{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return UserAuth{}, ErrInvalidTokenStructure
	}

	u := UserAuth{
		Id:   claims.UserId,
		Role: claims.UserRole,
	}

	return u, nil
}

type AuthClaims struct {
	UserId   string
	UserRole role.Role
	jwt.StandardClaims
}

type UserAuth struct {
	Id   string
	Role role.Role
}
