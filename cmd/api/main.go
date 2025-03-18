package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kevin-griley/api/data"
	"github.com/kevin-griley/api/db"
	"github.com/kevin-griley/api/docs"
	"github.com/kevin-griley/api/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Rest API
//	@description	API for the ULD Management System
//
// @BasePath					/
//
// @securityDefinitions.apikey	Bearer Authentication
// @tokenUrl http://localhost:3000/login
// @in							header
// @name						Authorization
// @description
func main() {

	listenAddress := ":3000"
	docs.SwaggerInfo.Host = "localhost:3000"

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// Serve the Swagger API documentation
	mux.HandleFunc("GET /docs/", httpSwagger.WrapHandler)

	// Auth routes
	mux.HandleFunc("POST /login", wrapError(handlers.HandlePostLogin))

	// User routes
	mux.HandleFunc("GET /user/{id}", withJWTAuth(wrapError((handlers.HandleGetUser))))
	mux.HandleFunc("POST /user", wrapError(handlers.HandlePostUser))

	log.Println("Open Dev Docs", "http://localhost:3000/docs")
	http.ListenAndServe(listenAddress, mux)
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

func withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling withJWTAuth middleware")

		tokenStr := r.Header.Get("Authorization")
		token, err := validateJWT(tokenStr)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}

		userID, err := handlers.GetID(r)
		if err != nil {
			permissionDenied(w)
			return
		}
		user, err := data.GetUserByID(userID)
		if err != nil {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		sub, err := claims.GetSubject()
		if err != nil {
			permissionDenied(w)
			return
		}

		userIDParsed, err := uuid.Parse(sub)
		if err != nil {
			permissionDenied(w)
			return
		}

		if user.ID != userIDParsed {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenStr string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	handlers.WriteJSON(w, http.StatusUnauthorized, handlers.ApiError{
		Status:  http.StatusUnauthorized,
		Message: "Permission Denied",
	})
}
