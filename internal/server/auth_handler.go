package server

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (s *Server) adminLoginHandler(w http.ResponseWriter, r *http.Request) {
	s.loginHandler(w, r, 0) // 0 for admin
}

func (s *Server) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	s.loginHandler(w, r, 1) // 1 for user
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request, role int16) {
	ctx := r.Context()
	var req LoginRequest
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

	token, err := s.service.UsernameLogin(ctx, req.Username, req.Password, role)
	if err != nil {
		log.Printf("failed to login: %s\n", err.Error())
		// if errors.Is(err, constants.ErrUserWrongPassword) {
		// 	sendErrorResponse(w, http.StatusBadRequest, "wrong password")
		// 	return
		// }
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := LoginResponse{
		Token: token,
	}
	sendResponse(w, http.StatusOK, resp)
}

func (s *Server) adminRegisterHandler(w http.ResponseWriter, r *http.Request) {
	s.registerHandler(w, r, 0) // 0 for admin
}

func (s *Server) userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	s.registerHandler(w, r, 1) // 1 for user
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request, role int16) {
	ctx := r.Context()
	var req RegisterRequest
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
		Username: sql.NullString{
			String: req.Username,
			Valid:  true,
		},
		Role: role,
	}
	token, err := s.service.Register(ctx, user, req.Password, role)
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
		Token: token,
	}
	sendResponse(w, http.StatusCreated, resp)
}
