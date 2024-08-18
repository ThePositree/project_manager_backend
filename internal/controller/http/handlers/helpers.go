package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, status int, jsonMessage any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	bytes, err := json.Marshal(jsonMessage)
	if err != nil {
		return fmt.Errorf("marshaling message json: %w", err)
	}
	w.Write(bytes)
	return nil
}

type ResponseMessageDTO struct {
	Message string `json:"message"`
}