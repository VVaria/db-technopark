package middleware

import (
	"net/http"
)

func appJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, req)
	})
}