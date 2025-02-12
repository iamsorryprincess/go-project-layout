package middleware

import "net/http"

func WithHandler(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func WithHandlerFunc(handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) http.Handler {
	return WithHandler(handler, middlewares...)
}
