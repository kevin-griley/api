package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
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
// @description					A valid JWT token with Bearer prefix
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load embedded .env file:", err)
	}

	dbConn, err := db.Init()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close(dbConn)

	listenAddress := ":3000"
	docs.SwaggerInfo.Host = "localhost:3000"

	mux := http.NewServeMux()
	//////////////////////////
	// SWAGGER - Docs route
	//////////////////////////
	mux.HandleFunc("GET /docs/", httpSwagger.WrapHandler)
	//////////////////////////
	// AUTH - Login route
	//////////////////////////
	PostLoginHandler := handlers.HandleApiError(handlers.HandlePostLogin)
	PostLogin := middleware.Chain(
		PostLoginHandler,
		middleware.DBMiddleware(dbConn),
	)
	mux.HandleFunc("POST /login", PostLogin)
	//////////////////////////
	// USER - Get Current User Route
	//////////////////////////
	GetUserByKeyHandler := handlers.HandleApiError(handlers.HandleGetUserByKey)
	GetUserByKey := middleware.Chain(
		GetUserByKeyHandler,
		middleware.JwtAuthMiddleware,
		middleware.DBMiddleware(dbConn),
	)
	mux.HandleFunc("GET /user/me", GetUserByKey)
	//////////////////////////
	// USER - Create New User Route
	//////////////////////////
	PostUserHandler := handlers.HandleApiError(handlers.HandlePostUser)
	PostUser := middleware.Chain(
		PostUserHandler,
		middleware.DBMiddleware(dbConn),
	)
	mux.HandleFunc("POST /user", PostUser)
	//////////////////////////
	// USER - Update User Route
	//////////////////////////
	PatchUserHandler := handlers.HandleApiError(handlers.HandlePatchUser)
	PatchUser := middleware.Chain(
		PatchUserHandler,
		middleware.JwtAuthMiddleware,
		middleware.DBMiddleware(dbConn),
	)
	mux.HandleFunc("PATCH /user/me", PatchUser)

	slog.Info("Application", "Swagger Docs Url", "http://localhost:3000/docs")

	loggingHandler := middleware.LoggingMiddleware(mux.ServeHTTP)

	http.ListenAndServe(listenAddress, loggingHandler)
}
