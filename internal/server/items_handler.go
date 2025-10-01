package server

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) createItemHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateItemRequest
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

	merchantIDStr := r.PathValue("merchantID")
	if len(merchantIDStr) == 0 {
		sendErrorResponse(w, http.StatusBadRequest, "product id cannot be empty")
	}
	merchantID, err := strconv.ParseInt(merchantIDStr, 10, 64)
	if err != nil || merchantID <= 0 {
		sendErrorResponse(w, http.StatusNotFound, "merchant not found")
		return
	}

	if !constants.IsValidProductCategory(req.ProductCategory) {
		sendErrorResponse(w, http.StatusBadRequest, "invalid product category")
		return
	}

	if req.ImageURL == "" || !(strings.HasPrefix(req.ImageURL, "http://") ||
		strings.HasPrefix(req.ImageURL, "https://")) {
		sendErrorResponse(w, http.StatusBadRequest, "invalid imageUrl")
		return
	}

	newItem := model.Item{
		MerchantID: merchantID,
		Name:       req.Name,
		Category:   req.ProductCategory,
		Price:      req.Price,
		ImageURL:   req.ImageURL,
	}

	itemID, err := s.service.CreateItems(ctx, newItem)
	if err != nil {
		if errors.Is(err, constants.ErrMerchantNotFound) {
			sendErrorResponse(w, http.StatusNotFound, "merchant not found")
			return
		}
		log.Printf("failed to create new item: %s\n", err.Error())
		sendErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := CreateItemResponse{
		itemID: strconv.Itoa(int(itemID)),
	}

	sendResponse(w, http.StatusCreated, response)
}
