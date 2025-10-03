package server

import (
	"PattyWagon/internal/constants"
	"encoding/json"
	"errors"
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

func handleServiceError(w http.ResponseWriter, err error, customMappings map[error]string) {
	for targetErr, customMsg := range customMappings {
		if errors.Is(err, targetErr) {
			switch {
			case errors.Is(err, constants.ErrUserNotFound):
				sendErrorResponse(w, http.StatusNotFound, customMsg)
				return
			case errors.Is(err, constants.ErrUserWrongPassword):
				sendErrorResponse(w, http.StatusBadRequest, customMsg)
				return
			case errors.Is(err, constants.ErrDuplicate):
				sendErrorResponse(w, http.StatusConflict, customMsg)
				return
			case errors.Is(err, constants.ErrDuplicateEmail):
				sendErrorResponse(w, http.StatusConflict, customMsg)
				return
			case errors.Is(err, constants.ErrDuplicatePhoneNum):
				sendErrorResponse(w, http.StatusConflict, customMsg)
				return
			}
		}
	}

	handleDefaultError(w, err)
}

func handleDefaultError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, constants.ErrMaximumFileSize):
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	case errors.Is(err, constants.ErrInvalidFileType):
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	case errors.Is(err, constants.ErrUserNotFound):
		sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	case errors.Is(err, constants.ErrUserWrongPassword):
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	case errors.Is(err, constants.ErrDuplicate):
		sendErrorResponse(w, http.StatusConflict, err.Error())
		return
	case errors.Is(err, constants.ErrDuplicateEmail):
		sendErrorResponse(w, http.StatusConflict, err.Error())
		return
	case errors.Is(err, constants.ErrDuplicatePhoneNum):
		sendErrorResponse(w, http.StatusConflict, err.Error())
		return
	default:
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

type LoginResponse struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
	Token string `json:"token"`
}

type RegisterResponse struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
	Token string `json:"token"`
}

type FileUploadData struct {
	ImageUrl string `json:"imageUrl"`
}

type FileUploadResponse struct {
	Message string         `json:"message"`
	Data    FileUploadData `json:"data"`
}
