package middleware

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(final http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		final = middlewares[i](final)
	}
	return final
}
