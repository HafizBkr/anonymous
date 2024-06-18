// middlewares/middlewares.go

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
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := m.jwt.Decode(token)
        if err != nil {
            m.logger.Error(fmt.Sprintf("Error decoding token: %s", err))
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        userId, ok := claims["id"].(string)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        user, err := m.users.GetUserDataByID(userId)
        if err != nil {
            m.logger.Error(fmt.Sprintf("Error fetching user: %s", err))
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if !user.Active {
            http.Error(w, "User inactive", http.StatusForbidden)
            return
        }

        ctx := context.WithValue(r.Context(), ContextKeyUser, user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
