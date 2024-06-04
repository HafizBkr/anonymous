package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"anonymous/auth"
	"anonymous/types"
)

// ContextKey est la clé pour stocker les valeurs dans le contexte de la requête
type ContextKey string

const ContextKeyUser ContextKey = "user"

// AuthMiddleware est un middleware qui vérifie l'authentification de l'utilisateur
type AuthMiddleware struct {
	users  auth.UserRepo
	jwt    types.JWTProvider
	logger types.Logger
}

// NewAuthMiddleware crée une nouvelle instance de AuthMiddleware
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

// MiddlewareHandler est un middleware qui vérifie l'authentification de l'utilisateur via le jeton JWT
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

		// Stocke l'utilisateur dans le contexte de la requête
		ctx := context.WithValue(r.Context(), ContextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
