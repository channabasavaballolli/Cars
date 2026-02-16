package middleware

import (
	"context"
	"net/http"
	"strings"

	"car-service/utils"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey = contextKey("userID")

// AuthMiddleware validates the JWT token in the Authorization header
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No token, proceed without user context (public access or handled by resolver)
			next.ServeHTTP(w, r)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"errors": [{"message": "Invalid Authorization header format"}]}`))
			return
		}

		tokenString := bearerToken[1]
		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"errors": [{"message": "Invalid or expired token"}]}`))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"errors": [{"message": "Invalid token claims"}]}`))
			return
		}

		// Extract user_id safely
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"errors": [{"message": "Invalid user ID in token"}]}`))
			return
		}
		userID := int(userIDFloat)

		// Add userID to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
