package utils

import "net/http"

// Middleware is a alias for func(http.Handler) http.Handler
type Middleware func(http.Handler) http.Handler

// PipelineHandlers is a util function for helping write a pipeline of http.Handler
func PipelineHandlers(middlewares []Middleware,
	handler http.Handler) http.Handler {

	for i := len(middlewares); i > 0; i-- {
		handler = middlewares[i-1](handler)
	}

	return handler
}

// PipelineHandlerFuncs is a util function for helping write a pipeline of http.HandlerFunc
func PipelineHandlerFuncs(middlewares []Middleware,
	handlerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

	var handler http.Handler = http.HandlerFunc(handlerFunc)

	for i := len(middlewares); i > 0; i-- {
		handler = middlewares[i-1](handler)
	}

	return handler.ServeHTTP
}
