package utils

import "net/http"

func PipelineRoutes(middlewares []func(http.Handler) http.Handler,
	handler http.Handler) http.Handler {

	for i := len(middlewares); i > 0; i-- {
		handler = middlewares[i-1](handler)
	}

	return handler
}
