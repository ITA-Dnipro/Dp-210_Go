package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware/auth"
)

type contextKey string

var (
	userIdKey   = contextKey("userId")
	userRoleKey = contextKey("userRole")
)

func (*Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, ok := tokenfromRequest(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		}

		claims, err := auth.ValidateToken(t)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		ctx := NewContext(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewContext(ctx context.Context, claims *auth.AuthClaims) context.Context {
	ctx = context.WithValue(ctx, userIdKey, claims.UserId)
	return context.WithValue(ctx, userRoleKey, claims.UserRole)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIdKey).(string)
	return id, ok
}
func UserRoleFromContext(ctx context.Context) (role.Role, bool) {
	r, ok := ctx.Value(userRoleKey).(role.Role)
	return r, ok
}

func tokenfromRequest(r *http.Request) (auth.JwtToken, bool) {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		return auth.JwtToken(""), false
	}

	jwtToken := auth.JwtToken(authHeader[1])
	return jwtToken, true
}
