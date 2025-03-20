package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MazenBahie/go_http_server/models"
	"github.com/gorilla/mux"
)

func HandleAddCreditCard(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cardDetails := models.CreditCard{}
		json.NewDecoder(r.Body).Decode(&cardDetails)

		err := Val.Struct(cardDetails)
		if err != nil {
			ResponseWithError(w, http.StatusBadRequest, "Validation error "+err.Error())
			return
		}

		createdCard := models.CreditCard{}

		query := `
            INSERT INTO credit_cards (name, card_number, user_id, expiration_date, cvv)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id, name, card_number, user_id, expiration_date, cvv
        `
		err = db.QueryRow(query, cardDetails.Name, cardDetails.CardNumber, cardDetails.UserID, cardDetails.ExpirationDate, cardDetails.Cvv).Scan(
			&createdCard.ID, &createdCard.Name, &createdCard.CardNumber, &createdCard.UserID, &createdCard.ExpirationDate, &createdCard.Cvv,
		)

		if err != nil {
			ResponseWithError(w, http.StatusInternalServerError, "Error adding credit card")
			return
		}

		ResponseJson(w, http.StatusCreated, createdCard)
	}
}

func HandleDeleteCreditCard(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["card_id"]
		if id == "" {
			ResponseWithError(w, http.StatusBadRequest, "card_id is required")
			return
		}

		castedId, err := strconv.Atoi(id)
		if err != nil {
			ResponseWithError(w, http.StatusBadRequest, "card_id must be a number")
			return
		}

		query := `
            DELETE FROM credit_cards
			WHERE id = $1
        `
		_, err = db.Exec(query, castedId)

		if err != nil {
			ResponseWithError(w, http.StatusInternalServerError, "Error deleting credit card: "+err.Error())
			return
		}

		ResponseJson(w, http.StatusOK, "")
	}
}
