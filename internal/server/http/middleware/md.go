package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	"go.uber.org/zap"
)

// UsersRepository represent user repository.
type UsersRepository interface {
	GetByID(ctx context.Context, id string) (entity.User, error)
}

type Middleware struct {
	Logger *zap.Logger
	UR     UsersRepository
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

func (m *Middleware) RoleOnly(roles ...role.Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			id, ok := FromContext(ctx)
			if ok {
				u, err := m.UR.GetByID(ctx, id)
				if err == nil && role.IsAllowedRole(role.Role(u.PermissionRole), roles) {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		})
	}
}
