package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/MazenBahie/go_http_server/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type signUpResult struct {
	UserName string `json:"username"`
	Role     string `json:"role"`
}

func HandleSignUp(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.User{}

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			JsonParsingErrorBadRequest(w, err.Error())
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			ResponseWithError(w, http.StatusInternalServerError, "Error hashing password "+err.Error())
			return
		}
		user.Password = string(hashedPassword)

		res, err := db.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", user.Username, user.Password, user.Role)
		if err != nil || res == nil {
			ResponseWithError(w, http.StatusInternalServerError, "Error creating user "+err.Error())
			return
		}

		ResponseJson(w, http.StatusCreated, signUpResult{
			UserName: user.Username,
			Role:     user.Role,
		})
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func HandleLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginCreds LoginRequest
		json.NewDecoder(r.Body).Decode(&loginCreds)

		var storedUser models.User
		err := db.QueryRow("SELECT id, username, password, role FROM users WHERE username = $1",
			loginCreds.Username).Scan(&storedUser.ID, &storedUser.Username, &storedUser.Password, &storedUser.Role)

		if err != nil {
			ResponseWithError(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}

		// Compare passwords
		err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(loginCreds.Password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(5 * time.Hour)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": storedUser.Username,
			"role":     storedUser.Role,
			"exp":      expirationTime.Unix(),
		})

		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			ResponseWithError(w, http.StatusInternalServerError, "Could not generate token: "+err.Error())
			return
		}

		w.Header().Set("Authorization", "Bearer "+tokenString)
		ResponseJson(w, http.StatusOK, LoginResponse{Token: tokenString})
	}
}
