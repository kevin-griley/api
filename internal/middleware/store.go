package middleware

import (
	"net/http"

	"github.com/kevin-griley/api/internal/data"
)

func StoreMiddleware(store *data.Store) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := data.WithStore(r.Context(), store)
			next(w, r.WithContext(ctx))

		}
	}
}
