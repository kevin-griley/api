package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kevin-griley/api/internal/data"
	"github.com/kevin-griley/api/internal/db"
	"github.com/kevin-griley/api/internal/middleware"
)

func TestLogin(t *testing.T) {

	dbConn, err := db.Init()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close(dbConn)

	store := data.NewStore(dbConn)

	validEmail := "Kevin"
	validPassword := "Kevin"
	wrongPassword := "Kev"

	testCases := []struct {
		name           string
		loginPayload   PostAuthRequest
		expectedStatus int
	}{
		{
			name: "Valid Login",
			loginPayload: PostAuthRequest{
				Email:    validEmail,
				Password: validPassword,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid Password",
			loginPayload: PostAuthRequest{
				Email:    validEmail,
				Password: wrongPassword,
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Empty Credentials",
			loginPayload: PostAuthRequest{
				Email:    "",
				Password: "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if apiErr := HandlePostLogin(w, r); apiErr != nil {
			http.Error(w, apiErr.Message, apiErr.Status)
		}
	})

	finalHandler := middleware.Chain(handler, middleware.StoreMiddleware(store))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			reqBody, err := json.Marshal(tc.loginPayload)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			finalHandler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response body: %s",
					tc.expectedStatus, rr.Code, rr.Body.String())
			}

			if tc.expectedStatus == http.StatusOK {
				var resp PostAuthResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if resp.Token == "" {
					t.Errorf("Expected non-empty JWT token in response")
				}

				token, err := middleware.ValidateJWT(resp.Token)
				if err != nil {
					t.Fatalf("Failed to validate JWT token: %v", err)
				}

				_, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					t.Fatalf("Failed to parse claims from JWT token")
				}

			}
		})
	}
}
