package server

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"PattyWagon/internal/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

	merchantIDStr := r.PathValue("merchantId")
	if len(merchantIDStr) == 0 {
		sendErrorResponse(w, http.StatusBadRequest, "merchant id cannot be empty")
		return
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

	if req.Name == "" {
		sendErrorResponse(w, http.StatusBadRequest, "item name cannot be empty")
		return
	}

	if len(req.Name) < 2 || len(req.Name) > 30 {
		sendErrorResponse(w, http.StatusBadRequest, "invalid item name")
		return
	}

	if req.ImageURL == "" || !(strings.HasPrefix(req.ImageURL, "http://") ||
		strings.HasPrefix(req.ImageURL, "https://")) {
		sendErrorResponse(w, http.StatusBadRequest, "invalid imageUrl")
		return
	}

	if err := utils.ValidateFileExtensions(filepath.Base(req.ImageURL), constants.AllowedExtensions); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Price == 0 {
		sendErrorResponse(w, http.StatusBadRequest, "invalid price: price cannot be null")
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
		ItemID: strconv.Itoa(int(itemID)),
	}

	sendResponse(w, http.StatusCreated, response)
}

func (s *Server) getItemHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodGet {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	query := r.URL.Query()
	req := GetItemsRequest{
		ItemID:          query.Get("itemID"),
		Limit:           query.Get("limit"),
		Offset:          query.Get("offset"),
		Name:            query.Get("name"),
		ProductCategory: query.Get("productCategory"),
		CreatedAt:       query.Get("createdAt"),
	}

	itemIDint, _ := strconv.Atoi(req.ItemID)
	limitInt, _ := strconv.Atoi(req.Limit)
	offsetInt, _ := strconv.Atoi(req.Offset)

	paramsItem := model.FilterItem{
		ItemID:          int64(itemIDint),
		Limit:           limitInt,
		Offset:          offsetInt,
		Name:            req.Name,
		ProductCategory: req.ProductCategory,
		CreatedAt:       req.CreatedAt,
	}

	if paramsItem.Limit <= 0 {
		paramsItem.Limit = 5
	}
	if paramsItem.Offset < 0 {
		paramsItem.Offset = 0
	}

	if paramsItem.ProductCategory != "" {
		if !constants.IsValidProductCategory(paramsItem.ProductCategory) {
			sendResponse(w, http.StatusOK, []model.Merchant{})
			return
		}
	}

	if paramsItem.CreatedAt != "" {
		if !strings.EqualFold(paramsItem.CreatedAt, "asc") &&
			!strings.EqualFold(paramsItem.CreatedAt, "desc") {
			paramsItem.CreatedAt = ""
		}
	}

	items, err := s.service.GetItems(ctx, paramsItem)
	if err != nil {
		log.Printf("failed to get merchants: %s\n", err.Error())
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	var detailItems []DetailItem
	var meta Meta
	if len(items) > 0 {
		for _, item := range items {
			detail := DetailItem{
				ItemID:          strconv.Itoa(int(item.ID)),
				Name:            item.Name,
				ProductCategory: item.Category,
				Price:           item.Price,
				ImageURL:        item.ImageURL,
				CreatedAt:       item.CreatedAt.Format(time.RFC3339),
			}

			detailItems = append(detailItems, detail)
		}
	} else {
		detailItems = []DetailItem{}
	}

	meta = Meta{
		Limit:  paramsItem.Limit,
		Offset: paramsItem.Offset,
		Total:  len(detailItems),
	}

	res := GetItemsResponse{
		Data: detailItems,
		Meta: meta,
	}

	sendResponse(w, http.StatusOK, res)
}
