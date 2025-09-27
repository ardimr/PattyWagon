package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.HelloWorldHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("POST /v1/register/email", s.emailRegisterHandler)
	mux.HandleFunc("POST /v1/login/email", s.emailLoginHandler)
	mux.HandleFunc("POST /v1/file", s.fileUploadHandler)

	return s.contentMiddleware(s.authMiddleware(mux))
}
