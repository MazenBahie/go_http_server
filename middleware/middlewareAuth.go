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
		// secretKey := "NewSecretKey" // Hardcoded secret key for validation

		// Parse the token
		claims := jwt.MapClaims{} // ✅ Use non-pointer
		// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		// 	}

		// 	return []byte(os.Getenv("JWT_SECRET")), nil
		// })
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			log.Println("Error parsing token:", err)
			handlers.ResponseWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Ensure correct signing method
		if token.Method != jwt.SigningMethodHS256 {
			log.Println("Unexpected signing method:", token.Header["alg"])
			http.Error(w, "Invalid token algorithm", http.StatusUnauthorized)
			return
		}

		// Extract claims
		role, ok := claims["role"].(string)
		if !ok {
			handlers.ResponseWithError(w, http.StatusUnauthorized, "Role not found in token")
			return
		}

		// Debugging logs
		log.Println("Decoded Role:", role)
		log.Println("JWT_SECRET used for validation:", os.Getenv("JWT_SECRET"))

		// Authorization check for admin-only routes
		if adminOnly && role != "admin" {
			handlers.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	}
}

// func AuthMiddleware(next http.Handler, adminOnly bool) http.HandlerFunc {
// 	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			handlers.ResponseWithError(w, http.StatusUnauthorized, "Missing token or not logged in")
// 			return
// 		}

// 		tokenString := strings.Split(authHeader, " ")[1]
// 		// key := os.Getenv("JWT_SECRET")

// 		claims := jwt.MapClaims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			// return []byte(key), nil
// 			return []byte("NewSecretKey"), nil
// 		})

// 		if err != nil {
// 			handlers.ResponseWithError(w, http.StatusUnauthorized, "Invalid token")
// 			return
// 		}

// 		if token.Method != jwt.SigningMethodHS256 {
// 			log.Println("Unexpected signing method:", token.Header["alg"])
// 			http.Error(w, "Invalid token algorithm", http.StatusUnauthorized)
// 			return
// 		}

// 		extractedClaims, ok := token.Claims.(jwt.MapClaims)
// 		if !ok {
// 			handlers.ResponseWithError(w, http.StatusUnauthorized, "Invalid token claims")
// 			return
// 		}

// 		role, ok := extractedClaims["role"].(string)
// 		if !ok {
// 			handlers.ResponseWithError(w, http.StatusUnauthorized, "Role not found in token")
// 			return
// 		}

// 		if adminOnly && role != "admin" {
// 			handlers.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized")
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 		// })
// 	}

// }
