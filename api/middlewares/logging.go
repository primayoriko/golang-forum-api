package middlewares

import "net/http"

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logging stuff

		next.ServeHTTP(w, r)
	})
}
