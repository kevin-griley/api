package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kevin-griley/api/internal/data"
	"github.com/kevin-griley/api/internal/db"
	"github.com/kevin-griley/api/internal/middleware"
)

func BaseLine(handlerFunc ApiFunc) (http.HandlerFunc, *jwt.Token, error) {

	dbConn, err := db.Init()
	if err != nil {
		return nil, nil, err
	}

	store := data.NewStore(dbConn)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if apiErr := handlerFunc(w, r); apiErr != nil {
			http.Error(w, apiErr.Message, apiErr.Status)
		}
	})

	loginPayload := PostAuthRequest{
		Email:    "Kevin",
		Password: "Kevin",
	}

	reqBody, err := json.Marshal(loginPayload)
	if err != nil {
		return nil, nil, err
	}

	loginHandler := middleware.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if apiErr := HandlePostLogin(w, r); apiErr != nil {
			http.Error(w, apiErr.Message, apiErr.Status)
		}
	}), middleware.StoreMiddleware(store))

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	loginHandler.ServeHTTP(rr, req)

	resp := new(PostAuthResponse)
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	token, err := middleware.ValidateJWT(resp.Token)
	if err != nil {
		return nil, token, fmt.Errorf("failed to validate JWT: %v", err)
	}

	handler = middleware.Chain(
		handler,
		middleware.JwtAuthMiddleware,
		middleware.StoreMiddleware(store),
	)

	return handler, token, nil

}
