package middleware

import "net/http"

// Middleware represents a middleware function
type Middleware func(http.Handler) http.Handler

// Chain applies multiple middleware functions to a handler
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// ApplyMiddleware is an alias for Chain for backward compatibility
func ApplyMiddleware(handler http.Handler, middlewares ...Middleware) http.Handler {
	return Chain(handler, middlewares...)
}
