package middleware

import (
	"context"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/auth"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/entity"
	"net/http"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	"go.uber.org/zap"
)

type UserUsecases interface {
	GetRoleByID(ctx context.Context, id string) (role.Role, error)
}

type Auth interface {
	ValidateToken(t auth.JwtToken) (auth.UserAuth, error)
}

type Middleware struct {
	Logger *zap.Logger
	UserUC UserUsecases
	Auth   Auth
}

func (m *Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		m.Logger.Info("incoming request", zap.String("URI", r.RequestURI))
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		start := time.Now()
		next.ServeHTTP(w, r)

		m.Logger.Info("request finished",
			zap.String("took", time.Since(start).String()),
		)
	})
}

func (m *Middleware) RoleOnly(roles ...entity.Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			u, ok := UserFromContext(ctx)
			if ok && entity.IsAllowedRole(u.Role, roles) {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		})
	}
}
