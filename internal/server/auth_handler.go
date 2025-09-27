package server

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/model"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (s *Server) emailLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req EmailLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("invalid login request")
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = s.validator.Struct(req)
	if err != nil {
		log.Printf("invalid login request: %s\n", err.Error())
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, phone, err := s.service.EmailLogin(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("failed to login: %s\n", err.Error())
		if errors.Is(err, constants.ErrUserNotFound) {
			sendErrorResponse(w, http.StatusNotFound, fmt.Sprintf("email %s not found", req.Email))
			return
		}
		if errors.Is(err, constants.ErrUserWrongPassword) {
			sendErrorResponse(w, http.StatusBadRequest, "wrong password")
			return
		}
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := LoginResponse{
		Email: req.Email,
		Phone: phone,
		Token: token,
	}
	sendResponse(w, http.StatusOK, resp)
}

func (s *Server) emailRegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req EmailRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("invalid login request")
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = s.validator.Struct(req)
	if err != nil {
		log.Printf("invalid register request: %s\n", err.Error())
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user := model.User{
		Email: sql.NullString{
			String: req.Email,
			Valid:  true,
		},
	}
	token, err := s.service.Register(ctx, user, req.Password)
	if err != nil {
		log.Printf("failed to register: %s\n", err.Error())
		if errors.Is(err, constants.ErrDuplicate) {
			sendErrorResponse(w, http.StatusConflict, fmt.Sprintf("email %s already exists", user.Email.String))
			return
		}
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := RegisterResponse{
		Email: req.Email,
		Token: token,
	}
	sendResponse(w, http.StatusCreated, resp)
}
