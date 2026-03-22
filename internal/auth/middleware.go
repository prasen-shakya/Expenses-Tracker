package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/prasen-shakya/todo/internal/users"
)

type contextKey string

const userContextKey contextKey = "user"

func RequireAuth(authService *Service, usersRepo *users.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Missing authorization header")
				return
			}

			tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
			if !ok || tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Invalid authorization header")
				return
			}

			userId, err := authService.VerifyJwtToken(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Invalid token")
				return
			}

			user, err := usersRepo.GetById(r.Context(), userId)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Invalid token")
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
