package server

import (
	"PattyWagon/logger"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.HelloWorldHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/admin/register", s.adminRegisterHandler)
	mux.HandleFunc("/admin/login", s.adminLoginHandler)
	mux.HandleFunc("/user/register", s.userRegisterHandler)
	mux.HandleFunc("/user/login", s.userLoginHandler)

	mux.HandleFunc("POST /v1/file", s.fileUploadHandler)

	return logger.LoggingMiddleware(s.contentMiddleware(s.authMiddleware(mux)))
}
