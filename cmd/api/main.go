package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/kevin-griley/api/db"
	docs "github.com/kevin-griley/api/docs"
	"github.com/kevin-griley/api/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			API
//	@description	API for Tofflemire Freight Services
//
// @BasePath					/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {

	listenAddress := ":3000"
	docs.SwaggerInfo.Host = "localhost:3000"

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /api/", httpSwagger.WrapHandler)

	router.HandleFunc("POST /login", wrapError(handlers.HandlePostLogin))
	router.HandleFunc("POST /register", wrapError(handlers.HandlePostUser))

	log.Println("Starting server on", listenAddress)

	http.ListenAndServe(listenAddress, router)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) *handlers.ApiError

func wrapError(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			slog.Error("API Error", "status", err.Status, "error", err.Error)
			handlers.WriteJSON(w, err.Status, err)
		}
	}
}
