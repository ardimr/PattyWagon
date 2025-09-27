package server

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/utils"
	"PattyWagon/logger"
	"PattyWagon/observability"
	"encoding/json"
	"net/http"
)

func (s *Server) GetNearbyMerchants(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.Tracer.Start(r.Context(), "handler.get_nearby_merchants")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

	if r.Method != http.MethodGet {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	coordinate := r.PathValue("coordinate")
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

	log.Printf("lat: %f, long: %f", lat, lng)

}

func (s *Server) EstimateOrderPrice(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.Tracer.Start(r.Context(), "handler.estimate_order_price")
	defer span.End()

	var req OrderEstimationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = s.validator.Struct(req)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	orderEstimationRequest := req.ToModel()

	estimationPrice, err := s.service.EstimateOrderPrice(ctx, orderEstimationRequest)
	if err != nil {
		switch err {
		case constants.ErrMerchantNotFound:
			sendErrorResponse(w, http.StatusNotFound, err.Error())
		case constants.ErrInvalidStartingPoint:
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		case constants.ErrMerchantTooFar:
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		default:
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	sendResponse(w, http.StatusOK, NewEstimationPriceResponse(estimationPrice))
}
