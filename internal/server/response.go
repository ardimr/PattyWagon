package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

func sendResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if body == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, error string) {
	resp := ErrorResponse{
		Error: error,
	}

	sendResponse(w, statusCode, resp)
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type FileUploadResponse struct {
	FileID           string `json:"fileId"`
	FileUri          string `json:"fileUri"`
	FileThumbnailUri string `json:"fileThumbnailUri"`
}
