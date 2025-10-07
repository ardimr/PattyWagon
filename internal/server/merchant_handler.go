package server

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"PattyWagon/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

	if err := utils.ValidateFileExtensions(filepath.Base(paramsCreateMerchant.ImageURL), constants.AllowedExtensions); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
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

func (s *Server) getMerchantHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodGet {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	query := r.URL.Query()
	req := GetMerchantRequest{
		MerchantID:       query.Get("merchantId"),
		Limit:            query.Get("limit"),
		Offset:           query.Get("offset"),
		Name:             query.Get("name"),
		MerchantCategory: query.Get("merchantCategory"),
		CreatedAt:        query.Get("createdAt"),
	}

	merchantIDint, _ := strconv.Atoi(req.MerchantID)
	limitInt, _ := strconv.Atoi(req.Limit)
	offsetInt, _ := strconv.Atoi(req.Offset)

	paramsMerchant := model.FilterMerchant{
		MerchantID:       int64(merchantIDint),
		Limit:            limitInt,
		Offset:           offsetInt,
		Name:             req.Name,
		MerchantCategory: req.MerchantCategory,
		CreatedAt:        req.CreatedAt,
	}

	if paramsMerchant.Limit <= 0 {
		paramsMerchant.Limit = 5
	}
	if paramsMerchant.Offset < 0 {
		paramsMerchant.Offset = 0
	}

	if paramsMerchant.MerchantCategory != "" {
		if !constants.IsValidMerchantCategory(paramsMerchant.MerchantCategory) {
			sendResponse(w, http.StatusOK, []model.Merchant{})
			return
		}
	}

	if paramsMerchant.CreatedAt != "" {
		if !strings.EqualFold(paramsMerchant.CreatedAt, "asc") &&
			!strings.EqualFold(paramsMerchant.CreatedAt, "desc") {
			paramsMerchant.CreatedAt = ""
		}
	}

	merchants, err := s.service.GetMerchants(ctx, paramsMerchant)
	if err != nil {
		log.Printf("failed to get merchants: %s\n", err.Error())
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	var detailMerchants []DetailMerchant
	var meta Meta
	if len(merchants) > 0 {
		for _, merchant := range merchants {
			detailMerchants = append(detailMerchants, DetailMerchant{
				MerchantID: strconv.Itoa(int(merchant.ID)),
				Name:       merchant.Name,
				Category:   utils.PointerValue(merchant.Category, ""),
				ImageURL:   merchant.ImageURL,
				Location: DetailLocation{
					Latitude:  merchant.Latitude,
					Longitude: merchant.Longitude,
				},
				CreatedAt: merchant.CreatedAt.Format(time.RFC3339),
			})
		}
	} else {
		detailMerchants = []DetailMerchant{}
	}

	meta = Meta{
		Limit:  paramsMerchant.Limit,
		Offset: paramsMerchant.Offset,
		Total:  len(detailMerchants),
	}

	res := GetMerchantResponse{
		Data: detailMerchants,
		Meta: meta,
	}

	sendResponse(w, http.StatusOK, res)
}
