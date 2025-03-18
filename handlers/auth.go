package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/MazenBahie/go_http_server/models"
	"golang.org/x/crypto/bcrypt"
)

type signUpResult struct {
	UserName string `json:"username"`
	Role     string `json:"role"`
}

var jwtKey string = os.Getenv("JWT_SECRET")

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

		// Insert into DB
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
