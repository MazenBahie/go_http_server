package handlers

import (
	"net/http"
)

func HandleError(w http.ResponseWriter, r *http.Request) {
	ResponseWithError(w, http.StatusInternalServerError, "Error occurred")
}
