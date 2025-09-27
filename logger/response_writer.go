package logger

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = append(rw.body, b...)
	return rw.ResponseWriter.Write(b)
}
