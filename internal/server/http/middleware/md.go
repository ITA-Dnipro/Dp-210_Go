package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"go.uber.org/zap"
)

type UserUsecases interface {
	GetRoleByID(ctx context.Context, id string) (entity.Role, error)
}

type Middleware struct {
	Logger *zap.Logger
	UserUC UserUsecases
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

func (m *Middleware) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RoleOnly(roles ...entity.Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			//id, ok := ctx.Value("id").(userId)
			ok := true
			if ok {
				//TODO replace get id from context.
				id := "text"
				role, err := m.UserUC.GetRoleByID(ctx, id)
				if err == nil && isAllowedRole(role, roles) {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		})
	}
}

func isAllowedRole(r entity.Role, allowedRoles []entity.Role) bool {
	for _, ar := range allowedRoles {
		if r == ar {
			return true
		}
	}
	return false
}
