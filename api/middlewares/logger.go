package middlewares

import "net/http"

// Log is a middleware method for writing log for every passed request
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logging stuff

		next.ServeHTTP(w, r)
	})
}
