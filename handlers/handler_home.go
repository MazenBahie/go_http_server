package handlers

import (
	"fmt"
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello, World!")
}
