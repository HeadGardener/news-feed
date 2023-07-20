package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Message string `json:"message"`
}

func newErrResponse(w http.ResponseWriter, code int, errorMsg string, err error) {
	log.Printf("[ERROR] %s: %e", errorMsg, err)
	newResponse(w, code, response{
		Message: errorMsg,
	})
}

func newResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
