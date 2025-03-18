package main

import (
	"log/slog"
	"net/http"

	"github.com/kevin-griley/api/docs"
	"github.com/kevin-griley/api/internal/db"
	"github.com/kevin-griley/api/internal/handlers"
	"github.com/kevin-griley/api/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Rest API
//	@description	API for the ULD Management System
//
// @BasePath					/
// @securityDefinitions.apikey	Bearer Authentication
// @tokenUrl http://localhost:3000/login
// @in							header
// @name						Authorization
// @description					Please provide a valid JWT token with Bearer prefix
func main() {

	listenAddress := ":3000"
	docs.SwaggerInfo.Host = "localhost:3000"

	if err := db.Init(); err != nil {
		slog.Error("Application", "Database Init", err)
	}

	mux := http.NewServeMux()

	// Serve the Swagger API documentation
	mux.HandleFunc("GET /docs/", httpSwagger.WrapHandler)

	// Auth routes
	mux.HandleFunc("POST /login", handlers.HandleApiError(handlers.HandlePostLogin))

	// User routes
	getUserByKeyHandler := handlers.HandleApiError(handlers.HandleGetUserByKey)
	getUserById := middleware.Chain(
		getUserByKeyHandler,
		middleware.JwtAuthMiddleware,
	)

	mux.HandleFunc("GET /user/me", getUserById)
	mux.HandleFunc("POST /user", handlers.HandleApiError(handlers.HandlePostUser))

	slog.Info("Application", "Swagger Docs Url", "http://localhost:3000/docs")

	finalHandler := middleware.LoggingMiddleware(mux.ServeHTTP)
	http.ListenAndServe(listenAddress, finalHandler)
}
