package server

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/utils"
	"context"
	"net/http"
	"strings"
)

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	publicPaths := map[string]bool{
		"/health":            true,
		"/v1/login/email":    true,
		"/v1/login/phone":    true,
		"/v1/register/email": true,
		"/v1/register/phone": true,
		"/v1/file":           true,
		"/v1/purchase":       true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if publicPaths[path] || (path == "/v1/product" && r.Method == "GET") || strings.HasPrefix(path, "/v1/purchase") {
			next.ServeHTTP(w, r)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			sendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		userID, err := utils.ParseUserIDFromToken(tokenString)
		if err != nil {
			sendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), constants.UserIDCtxKey, userID)

		// Proceed with the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) contentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/v1/file" {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			ct := r.Header.Get("Content-Type")
			if !strings.EqualFold(ct, "application/json") {
				http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
				return
			}
		}

		// continue to the next handler
		next.ServeHTTP(w, r)
	})
}
