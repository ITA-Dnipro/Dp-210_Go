package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
)

type contextKey string
type JwtToken string

type User struct {
	Id   string    `userId`
	Role role.Role `userRole`
}

type UserToken struct {
	Token JwtToken `json:"token"`
}

var (
	KeyUser = contextKey("user")
)

func (md *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, ok := tokenfromRequest(r)
		_ = t
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		}

		jsoned, err := json.Marshal(UserToken{Token: t})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp, err := http.Post(md.AuthUrl, "application/json", bytes.NewReader(jsoned))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		code := resp.StatusCode
		if code != http.StatusOK {
			w.WriteHeader(code)
			w.Write([]byte(resp.Status))
			return
		}

		var user User
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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

func tokenfromRequest(r *http.Request) (JwtToken, bool) {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		return "", false
	}

	jwtToken := JwtToken(authHeader[1])
	return jwtToken, true
}

type ReqUser struct {
	Id   string
	Role role.Role
}
