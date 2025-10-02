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
		"/health":         true,
		"/admin/register": true,
		"/admin/login":    true,
		"/user/register":  true,
		"/user/login":     true,
	}

	userPaths := map[string]bool{
		"/merchants/nearby/{lat}/{long}": true, //user dynamic lat long
		"/users/estimate":                true,
		"/users/orders":                  true,
	}

	adminPaths := map[string]bool{
		"/image":                              true,
		"/admin/merchants":                    true,
		"/admin/merchants/{merchantID}/items": true, //user dynamic merchantID
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if publicPaths[path] {
			next.ServeHTTP(w, r)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			sendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		userID, role, err := utils.ParseUserIDandRoleFromToken(tokenString)
		if err != nil {
			sendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if role == 1 {
			for _, pattern := range adminPaths {
				if pathMatchesPattern(pattern, path) {
					sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized User")
					return
				}
			}
		} else if role == 0 {
			for _, pattern := range userPaths {
				if pathMatchesPattern(pattern, path) {
					sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized Admin")
					return
				}
			}
		}

		ctx := context.WithValue(r.Context(), constants.UserIDCtxKey, userID)

		// Proceed with the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) contentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/image" {
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

func pathMatchesPattern(pattern, actual string) bool {
	patternParts := strings.Split(pattern, "/")
	actualParts := strings.Split(actual, "/")
	if len(patternParts) != len(actualParts) {
		return false
	}

	for i := range patternParts {
		if strings.HasPrefix(patternParts[i], "{") && strings.HasSuffix(patternParts[i], "}") {
			// It's a dynamic segment like {merchantID}, allow it to match anything
			continue
		}
		if patternParts[i] != actualParts[i] {
			return false
		}
	}

	return true
}
