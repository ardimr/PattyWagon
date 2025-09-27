package server

import (
	"PattyWagon/logger"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.HelloWorldHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("POST /v1/register/email", s.emailRegisterHandler)
	mux.HandleFunc("POST /v1/login/email", s.emailLoginHandler)
	mux.HandleFunc("POST /v1/file", s.fileUploadHandler)

	// Purchase
	mux.HandleFunc("GET /v1/merchants/nearby/{coordinate}", s.GetNearbyMerchants)
	mux.HandleFunc("POST /v1/users/estimate", s.EstimateOrderPrice)
	return logger.LoggingMiddleware(s.contentMiddleware(s.authMiddleware(mux)))
}
