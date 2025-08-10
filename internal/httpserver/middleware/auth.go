package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"work-tracker/internal/auth"
	"work-tracker/internal/config"
	"work-tracker/internal/store"
)

type contextKey string

const userIDKey contextKey = "uid"

func WithAuth(cfg config.Config, stores store.Stores) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			parts := strings.SplitN(authz, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			token := parts[1]
			claims, err := auth.ParseToken(token, cfg.JWTSecret)
			if err != nil || claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			exists, err := stores.Sessions.Exists(r.Context(), claims.ID)
			if err != nil || !exists {
				http.Error(w, "session expired", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDKey)
	if v == nil {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}
