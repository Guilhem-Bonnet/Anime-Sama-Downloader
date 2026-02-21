package httpjson

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("httpjson.Write: failed to encode response: %v", err)
	}
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	Write(w, status, ErrorResponse{Error: msg})
}
