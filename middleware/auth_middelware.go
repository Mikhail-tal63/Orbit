package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Mikhail-Tal63/Orbit/configs"
	"github.com/Mikhail-Tal63/Orbit/utils"
	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("missing authorization header"))
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid authorization format"))
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		secret := []byte(configs.Load().JWTSecret)
		userID, err := utils.VerifyJWT(secret, tokenString)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user id not found in context")
	}
	return userID, nil
}
