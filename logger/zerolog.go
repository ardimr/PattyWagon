package logger

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	if os.Getenv("ENV") == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log incoming request
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Msg("Request started")

		// Process next request
		wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
		next.ServeHTTP(wrapped, r.WithContext(r.Context()))

		// Log request completion
		logEvent := log.Info()
		if wrapped.statusCode >= 400 {
			logEvent = log.Error()
		}

		logEvent.Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", wrapped.statusCode)

		if wrapped.statusCode >= 400 && len(wrapped.body) > 0 {
			var errResp ErrorResponse

			if err := json.Unmarshal(wrapped.body, &errResp); err != nil {
				logEvent = logEvent.Str("error", errResp.Error)
			}
		}

		logEvent.Msg("Request completed")
	})
}
