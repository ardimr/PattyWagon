package server

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/utils"
	"PattyWagon/logger"
	"PattyWagon/observability"
	"net/http"
)

func (s *Server) FindNearbyMerchants(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerFromContext(r.Context())
	ctx, span := observability.Tracer.Start(r.Context(), "handler.get_nearby_merchants")
	defer span.End()
	log.Printf("------------finding nearby merchants ----------\n")

	if r.Method != http.MethodGet {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	coordinate := r.PathValue("coordinate")
	if coordinate == "" {
		sendErrorResponse(w, http.StatusBadRequest, "coordinate must not be empty")
		return
	}

	lat, lng, err := utils.ValidateAndExtractCoordinate(coordinate)
	if err != nil {
		switch err {
		case constants.ErrInvalidCoordinate:
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		default:
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	log.Printf("coordinate: %s | lat: %f lon: %f", coordinate, lat, lng)
	userLocation := LocationRequest{
		Lat:  lat,
		Long: lng,
	}

	query := r.URL.Query()
	searchParams := FindNearbyMerchantRequest{
		MerchantID:       query.Get("merchantId"),
		Limit:            query.Get("limit"),
		Offset:           query.Get("offset"),
		Name:             query.Get("name"),
		MerchantCategory: query.Get("merchantCategory"),
	}

	filter := searchParams.ToModel()
	merchants, err := s.service.FindNearbyMerchants(ctx, userLocation.ToModel(), filter)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
	meta := FindNearbyMerchantsResponseMeta{
		Limit:  filter.Limit,
		Offset: filter.Offset,
		Total:  len(merchants),
	}
	resp := NewFindNearbyMerchantsResponse(merchants, meta)
	sendResponse(w, http.StatusOK, resp)
}
