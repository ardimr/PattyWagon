package server

import (
	"PattyWagon/internal/model"
	"PattyWagon/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) createMerchantHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CreateMerchantRequest
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ct := r.Header.Get("Content-Type")
	if ct == "" || !strings.HasPrefix(ct, "application/json") {
		sendErrorResponse(w, http.StatusBadRequest, "invalid content type")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("invalid login request")
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := s.validator.Struct(req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "request doesn't pass validation")
		return
	}

	userID, ok := utils.GetUserIDFromCtx(ctx)
	if !ok {
		sendErrorResponse(w, http.StatusBadRequest, "Value not found or wrong type")
	}

	paramsCreateMerchant := model.Merchant{
		UserID:    userID,
		Name:      req.Name,
		Category:  utils.ToPointer(req.Category),
		ImageURL:  req.ImageURL,
		Latitude:  req.Location.Latitude,
		Longitude: req.Location.Longitude,
	}

	res, err := s.service.CreateMerchant(ctx, paramsCreateMerchant)
	if err != nil {
		log.Printf("failed to create new merchant: %s\n", err.Error())
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := CreateMerchantResponse{
		MerchantID: strconv.Itoa(int(res)),
	}

	sendResponse(w, http.StatusCreated, response)
}
