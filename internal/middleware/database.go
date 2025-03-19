package middleware

import (
	"database/sql"
	"net/http"

	"github.com/kevin-griley/api/internal/db"
)

func DBMiddleware(dbConn *sql.DB) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := db.WithDB(r.Context(), dbConn)
			next(w, r.WithContext(ctx))
		}
	}
}
