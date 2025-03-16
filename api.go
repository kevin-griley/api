package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	docs "github.com/kevin-griley/api/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	jwt "github.com/golang-jwt/jwt/v5"
)

type ApiConfig struct {
	listenAddress string
	store         Storage
}

func NewApiConfig() ApiConfig {

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	return ApiConfig{
		listenAddress: ":3000",
		store:         store,
	}
}

type APIServer struct {
	config ApiConfig
}

func NewAPIServer(config ApiConfig) *APIServer {
	return &APIServer{config}
}

//	@title			API
//	@description	API for Tofflemire Freight Services
//
// @BasePath					/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func (s *APIServer) Run() {

	// Docs
	docs.SwaggerInfo.Host = "localhost:3000"

	router := http.NewServeMux()

	router.HandleFunc("GET /api/", httpSwagger.WrapHandler)

	router.HandleFunc("POST /login", makeHandlerFunc(s.handleLogin))
	router.HandleFunc("POST /register", makeHandlerFunc(s.handleCreateUser))

	router.HandleFunc("GET /account", makeHandlerFunc(s.handleGet))
	router.HandleFunc("GET /account/{id}", withJWTAuth(makeHandlerFunc(s.handleGetByID), s.config.store))
	router.HandleFunc("POST /account", makeHandlerFunc(s.handleCreateAccount))
	router.HandleFunc("DELETE /account/{id}", makeHandlerFunc(s.handleDeleteAccount))

	log.Println("Starting server on", s.config.listenAddress)

	http.ListenAndServe(s.config.listenAddress, router)
}

type TokenResponse struct {
	Token string `json:"token"`
}

// @Summary			Login
// @Description		Login to the system
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			body	body		LoginRequest	true	"Login Request"
// @Success			200		{object}	TokenResponse	"Token Response"
// @Failure			400		{object} 	ApiError	"Bad Request"
// @Failure			401		{object} 	ApiError	"Unauthorized"
// @Router			/login	[post]
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) *ApiError {
	rBody := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(rBody); err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	user, err := s.config.store.GetUserByEmail(rBody.Email)
	if err != nil {
		return &ApiError{http.StatusUnauthorized, "invalid user or password"}
	}

	if !user.ValidPassword(rBody.Password) {
		return &ApiError{http.StatusUnauthorized, "invalid user or password"}
	}

	tokenString, err := createJWT(user)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}

	return WriteJSON(w, http.StatusOK, TokenResponse{Token: tokenString})
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) *ApiError {
	registerReq := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(registerReq); err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	user, err := NewUser(registerReq.Email, registerReq.Password)
	if err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	userID, err := s.config.store.CreateUser(user)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}

	user.ID = userID
	return WriteJSON(w, http.StatusOK, user)

}

func (s *APIServer) handleGet(w http.ResponseWriter, r *http.Request) *ApiError {
	accounts, err := s.config.store.GetAccounts()
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetByID(w http.ResponseWriter, r *http.Request) *ApiError {
	idStr, err := getID(r)
	if err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	account, err := s.config.store.GetAccountByID(idStr)
	if err != nil {
		return &ApiError{http.StatusNotFound, err.Error()}
	}
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) *ApiError {

	idStr, err := getID(r)
	if err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	if err := s.config.store.DeleteAccount(idStr); err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"deleted": idStr})
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) *ApiError {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	account := NewAccount(createAccountReq.Name)
	accountID, err := s.config.store.CreateAccount(account)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}

	account.ID = accountID

	return WriteJSON(w, http.StatusOK, account)

}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, "permission denied")

}

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		userID, err := getID(r)
		if err != nil {
			permissionDenied(w)
			return
		}
		account, err := s.GetAccountByID(userID)
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

		fmt.Println("sub: ", sub)
		fmt.Println("account.ID: ", account.ID)

		if account.ID != sub {
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

func WriteJSON(w http.ResponseWriter, status int, v interface{}) *ApiError {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}
	return nil
}

type apiFunc func(w http.ResponseWriter, r *http.Request) *ApiError

type ApiError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func makeHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			slog.Error("API Error", "status", err.Status, "error", err.Error)
			WriteJSON(w, err.Status, err)
		}
	}
}

func createJWT(user *User) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "mycartage",
		Subject:   user.ID,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))

}

func getID(r *http.Request) (string, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return idStr, fmt.Errorf("id is required")
	}
	return idStr, nil
}
