package httputil

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorBody struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {

		log.Printf("httputil: failed to encode JSON response: %v", err)
	}
}

func WriteError(w http.ResponseWriter, status int, message string) {
	var body errorBody
	body.Error.Message = message
	WriteJSON(w, status, body)
}
