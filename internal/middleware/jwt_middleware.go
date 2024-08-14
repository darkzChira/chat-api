package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

type JWTMiddleware struct {
	Secret []byte
}

func NewJWTMiddleware(secret []byte) *JWTMiddleware {
	return &JWTMiddleware{
		Secret: secret,
	}
}

func (m *JWTMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return m.Secret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		// Pass user ID from token to context
		userID := claims["user_id"].(string)
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
