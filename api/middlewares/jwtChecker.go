package middlewares

import "net/http"

// JwtCheck is a method for checking jwt token and passing the creds
func JwtCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logging stuff

		next.ServeHTTP(w, r)
	})
}
