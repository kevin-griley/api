package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type APIServer struct {
	listenAddress string
	store         Storage
}

func NewAPIServer(listenAddress string, store Storage) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store:         store,
	}
}

func (s *APIServer) Run() {
	router := http.NewServeMux()

	router.HandleFunc("POST /login", makeHandlerFunc(s.handleLogin))
	router.HandleFunc("POST /register", makeHandlerFunc(s.handleCreateUser))

	router.HandleFunc("GET /account", makeHandlerFunc(s.handleGet))
	router.HandleFunc("GET /account/{id}", withJWTAuth(makeHandlerFunc(s.handleGetByID), s.store))
	router.HandleFunc("POST /account", makeHandlerFunc(s.handleCreateAccount))
	router.HandleFunc("DELETE /account/{id}", makeHandlerFunc(s.handleDeleteAccount))

	log.Println("Starting server on", s.listenAddress)
	http.ListenAndServe(s.listenAddress, router)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) *ApiError {
	rBody := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(rBody); err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	user, err := s.store.GetUserByEmail(rBody.Email)
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

	return WriteJSON(w, http.StatusOK, map[string]string{"token": tokenString})
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

	userID, err := s.store.CreateUser(user)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}

	user.ID = userID
	return WriteJSON(w, http.StatusOK, user)

}

func (s *APIServer) handleGet(w http.ResponseWriter, r *http.Request) *ApiError {
	accounts, err := s.store.GetAccounts()
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

	account, err := s.store.GetAccountByID(idStr)
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

	if err := s.store.DeleteAccount(idStr); err != nil {
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
	accountID, err := s.store.CreateAccount(account)
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
