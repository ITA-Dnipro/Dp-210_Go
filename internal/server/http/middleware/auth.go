package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware/auth"
)

var UserIdContext = "userId"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		}

		jwtToken := authHeader[1]
		uId, err := auth.ValidateToken(jwtToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		ctx := context.WithValue(r.Context(), UserIdContext, uId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
