package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MazenBahie/go_http_server/handlers"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler, adminOnly bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			handlers.ResponseWithError(w, http.StatusUnauthorized, "Missing token or not logged in")
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]

		if tokenString == "" {
			handlers.ResponseWithError(w, http.StatusUnauthorized, "Missing token")
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			log.Println("Error parsing token:", err)
			handlers.ResponseWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		if token.Method != jwt.SigningMethodHS256 {
			handlers.ResponseWithError(w, http.StatusUnauthorized, "Invalid token algorithm")
			return
		}

		// Extract claims
		role, ok := claims["role"].(string)
		if !ok {
			handlers.ResponseWithError(w, http.StatusUnauthorized, "Role not found in token")
			return
		}

		// Authorization check for admin-only routes
		if adminOnly && role != "admin" {
			handlers.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	}
}
