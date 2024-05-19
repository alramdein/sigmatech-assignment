package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

func JWTMiddleware(secretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// Check if the token starts with "Bearer "
			tokenString := strings.TrimSpace(authHeader)
			if !strings.HasPrefix(tokenString, "Bearer ") {
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}

			// Extract the token
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			// Parse the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate the alg is what you expect
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, http.ErrAbortHandler
				}
				return secretKey, nil
			})

			// Check token validity
			if err != nil || !token.Valid {
				logrus.Error(err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Pass the execution to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
