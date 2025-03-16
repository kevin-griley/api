package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type APIServer struct {
	listenAddress string
	store 	   Storage
}

func NewAPIServer(listenAddress string, store Storage) *APIServer {
	return &APIServer{
		listenAddress: 	listenAddress,
		store: 			store,
	}
}


func (s *APIServer) Run() {
	router := http.NewServeMux()

	router.HandleFunc("GET /login", makeHTTPHandlerFunc(s.handleLogin))

	router.HandleFunc("GET /account", makeHTTPHandlerFunc(s.handleGet))
	router.HandleFunc("GET /account/{id}", withJWTAuth(makeHTTPHandlerFunc(s.handleGetByID), s.store))
	router.HandleFunc("POST /account", makeHTTPHandlerFunc(s.handleCreateAccount))
	router.HandleFunc("DELETE /account/{id}", makeHTTPHandlerFunc(s.handleDeleteAccount))


	log.Println("Starting server on", s.listenAddress)
	http.ListenAndServe(s.listenAddress, router)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (s *APIServer) handleGet(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetByID(w http.ResponseWriter, r *http.Request) error {
	idStr, err := getID(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByID(idStr)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	idStr, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(idStr); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"deleted": idStr})
}


func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.Name)
	userID, err := s.store.CreateAccount(account)
	if err != nil {
		return err
	}

	account.ID = userID
	tokenString, err := createJWT(account)
	if err != nil {
		return err
	}

	fmt.Println("JWT token: ", tokenString)

	return WriteJSON(w, http.StatusOK, createAccountReq)

}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, ApiError{"Permission Denied"})
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJteWNhcnRhZ2UiLCJzdWIiOiJhY2NvdW50IiwiZXhwIjoxNzQyMDk4NDA2LCJpYXQiOjE3NDIwMTIwMDZ9.jUss2yGnFTcSpOHpiPDVN5O5ieqCS6bWI1nh62m8-7w

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
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

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{err.Error()})
		}
	}
}


func createJWT(account *Account) (string, error) {

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Issuer: "mycartage",
		Subject: account.ID,

	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))


}

func getID (r *http.Request) (string, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return idStr, fmt.Errorf("id is required")
	}
	return idStr, nil
}