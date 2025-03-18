package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResponseJson(w http.ResponseWriter, code int, data interface{}) {
	dataAr, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dataAr)
}

type errorResponse struct {
	Message string `json:"message"`
}

func ResponseWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Internal server error 5xx: %v", msg)
	}

	ResponseJson(w, code, errorResponse{Message: msg})
}

func JsonParsingErrorBadRequest(w http.ResponseWriter, msg string) {
	ResponseWithError(w, http.StatusBadRequest, "Error parsing JSON: "+msg)
}
