package server

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/utils"
	"PattyWagon/logger"
	"PattyWagon/observability"
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
