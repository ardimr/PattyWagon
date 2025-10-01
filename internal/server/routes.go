package server

import (
	"PattyWagon/logger"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.HelloWorldHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("POST /admin/register", s.adminRegisterHandler)
	mux.HandleFunc("POST /admin/login", s.adminLoginHandler)
	mux.HandleFunc("POST /users/register", s.userRegisterHandler)
	mux.HandleFunc("POST /users/login", s.userLoginHandler)

	mux.HandleFunc("POST /v1/file", s.fileUploadHandler)
	mux.HandleFunc("POST /admin/merchant", s.createMerchantHandler)
	mux.HandleFunc("GET /admin/merchant", s.getMerchantHandler)
	mux.HandleFunc("POST /admin/merchants/{merchantId}/items", s.createItemHandler)

	return logger.LoggingMiddleware(s.contentMiddleware(s.authMiddleware(mux)))
}
