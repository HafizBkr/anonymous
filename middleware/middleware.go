package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"anonymous/auth"
	"anonymous/types"
)
type ContextKey string

const ContextKeyUser ContextKey = "user"

type AuthMiddleware struct {
	users  auth.UserRepo
	jwt    types.JWTProvider
	logger types.Logger
}

func NewAuthMiddleware(
	users auth.UserRepo,
	jwt types.JWTProvider,
	logger types.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		users:  users,
		jwt:    jwt,
		logger: logger,
	}
}
func (m *AuthMiddleware) MiddlewareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Non autorisé", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.jwt.Decode(token)
		if err != nil {
			m.logger.Error(fmt.Sprintf("Erreur lors du décodage du jeton : %s", err))
			http.Error(w, "Non autorisé", http.StatusUnauthorized)
			return
		}

		userId, ok := claims["id"].(string)
		if !ok {
			http.Error(w, "Non autorisé", http.StatusUnauthorized)
			return
		}

		user, err := m.users.GetUserDataByID(userId)
		if err != nil {
			m.logger.Error(fmt.Sprintf("Erreur lors de la récupération de l'utilisateur : %s", err))
			http.Error(w, "Non autorisé", http.StatusUnauthorized)
			return
		}

		if !user.Active {
			http.Error(w, "Utilisateur inactif", http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), ContextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

