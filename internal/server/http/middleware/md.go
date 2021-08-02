package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Middleware struct {
	Logger *zap.Logger
}

func (m *Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		m.Logger.Info("incoming request", zap.String("URI", r.RequestURI))
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		start := time.Now()
		next.ServeHTTP(w, r)

		m.Logger.Info("request finished",
			zap.String("took", time.Now().Sub(start).String()),
		)
	})
}
