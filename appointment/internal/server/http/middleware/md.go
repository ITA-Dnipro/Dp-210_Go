package middleware

import (
	"net/http"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/role"
	"go.uber.org/zap"
)

type Middleware struct {
	Logger *zap.Logger
}

func NewMiddleware(logger *zap.Logger) *Middleware {
	return &Middleware{Logger: logger}
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
			ur, ok := UserRoleFromContext(ctx)
			if ok && role.IsAllowedRole(role.Role(ur), roles) {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		})
	}
}
