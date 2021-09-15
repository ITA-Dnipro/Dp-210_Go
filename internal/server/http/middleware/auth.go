package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware/proto"
)

type contextKey string
type JwtToken string

type User struct {
	Id   string    `json:"userId"`
	Role role.Role `json:"userRole"`
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

		res, err := md.grpcClient.Validate(r.Context(), &proto.Token{Token: t})
		if err != nil {

			md.Logger.Warn(fmt.Sprintf("could not validate token via grpc %v", err))
			w.WriteHeader(500)
			return
		}

		if res.StatusCode != 200 {
			w.WriteHeader(int(res.StatusCode))
			return
		}

		ctx := contextWithUser(r.Context(), ReqUser{Id: res.UserId, Role: role.Role(res.UserRole)})
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

func tokenfromRequest(r *http.Request) (string, bool) {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		return "", false
	}

	jwtToken := authHeader[1]
	return jwtToken, true
}

type ReqUser struct {
	Id   string
	Role role.Role
}
