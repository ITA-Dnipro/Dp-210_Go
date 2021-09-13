package middleware

import (
	"net/http"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/role"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/service/auth"
	"go.uber.org/zap"
)

type UserUsecases interface {
	//	GetRoleByID(ctx context.Context, id string) (role.Role, error)
}

type Auth interface {
	ValidateToken(t auth.JwtToken) (auth.UserAuth, error)
}

type Middleware struct {
	Logger *zap.Logger
	UserUC UserUsecases
	Auth   Auth
}

func (md *Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		md.Logger.Info("incoming request", zap.String("URI", r.RequestURI))
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		start := time.Now()
		next.ServeHTTP(w, r)

		md.Logger.Info("request finished",
			zap.String("took", time.Since(start).String()),
		)
	})
}

func (md *Middleware) RoleOnly(roles ...role.Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			u, ok := UserFromContext(ctx)
			if ok && role.IsAllowedRole(u.Role, roles) {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		})
	}
}
