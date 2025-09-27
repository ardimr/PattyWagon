package server

import (
	"WeenieHut/internal/utils"
	"log"
	"net/http"
)

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetUserIDFromCtx(r.Context()))
	sendResponse(w, http.StatusOK, "Hello World!")
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, "OK")
}
