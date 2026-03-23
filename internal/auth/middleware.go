package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/prasen-shakya/todo/internal/respond"
	"github.com/prasen-shakya/todo/internal/users"
)

type contextKey string

const userContextKey contextKey = "user"

func RequireAuth(authService *Service, usersRepo *users.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respond.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Missing authorization header"})
				return
			}

			tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
			if !ok || tokenString == "" {
				respond.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header"})
				return
			}

			userId, err := authService.VerifyJwtToken(tokenString)
			if err != nil {
				respond.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
				return
			}

			user, err := usersRepo.GetById(r.Context(), userId)
			if err != nil {
				respond.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) (users.User, bool) {
	user, ok := ctx.Value(userContextKey).(users.User)
	return user, ok
}
