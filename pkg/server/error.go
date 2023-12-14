package server

import (
	"fmt"
	"net/http"
)

func SendInternalServerError(w http.ResponseWriter) {
	SendError(w, "internal server error", http.StatusInternalServerError)
}

func SendNotFound(w http.ResponseWriter) {
	SendError(w, "resource not found", http.StatusNotFound)
}

func SendError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"error": "%s"}`, message)
}
