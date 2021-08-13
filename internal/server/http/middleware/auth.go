package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware/auth"
)

type contextKey string

var (
	KeyUserId = contextKey("userId")
)

func (*Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, ok := tokenfromRequest(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		}

		uId, err := auth.ValidateToken(t)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		ctx := contextWithToken(r.Context(), uId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func contextWithToken(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, KeyUserId, userId)
}

func FromContext(ctx context.Context) (userId string, ok bool) {
	userId, ok = ctx.Value(KeyUserId).(string)
	return userId, ok
}

func tokenfromRequest(r *http.Request) (auth.JwtToken, bool) {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		return auth.JwtToken(""), false
	}

	jwtToken := auth.JwtToken(authHeader[1])
	return jwtToken, true
}
