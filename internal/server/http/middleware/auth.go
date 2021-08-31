package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/service/auth"
)

type contextKey string

var (
	KeyUser = contextKey("user")
)

func (md *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, ok := tokenfromRequest(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		}

		user, err := md.Auth.ValidateToken(t)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		ctx := contextWithUser(r.Context(), ReqUser{Id: user.Id, Role: user.Role})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func contextWithUser(ctx context.Context, user ReqUser) context.Context {
	return context.WithValue(ctx, KeyUser, user)
}

func UserFromContext(ctx context.Context) (user ReqUser, ok bool) {
	user, ok = ctx.Value(KeyUser).(ReqUser)
	return user, ok
}

func tokenfromRequest(r *http.Request) (auth.JwtToken, bool) {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		return auth.JwtToken(""), false
	}

	jwtToken := auth.JwtToken(authHeader[1])
	return jwtToken, true
}

type ReqUser struct {
	Id   string
	Role role.Role
}
