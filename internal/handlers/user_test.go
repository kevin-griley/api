package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kevin-griley/api/internal/data"
)

func TestUsers(t *testing.T) {

	number := rand.Int()
	newUserName := fmt.Sprintf("Kevin Griley %d", number)

	testCases := []struct {
		name           string
		method         string
		path           string
		updatePayload  PatchUserRequest
		expectedStatus int
	}{

		{
			name:   "Valid Update",
			method: http.MethodPatch,
			path:   "/user/me",
			updatePayload: PatchUserRequest{
				UserName: newUserName,
				Password: "Kevin",
			},
			expectedStatus: http.StatusOK,
		},

		{
			name:           "Get User",
			method:         http.MethodGet,
			path:           "/user/me",
			expectedStatus: http.StatusOK,
		},
	}

	finalHandler, token, err := BaseLine(HandlePatchUser)

	if err != nil {
		t.Fatalf("Failed to create baseline: %v", err)
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			reqBody, err := json.Marshal(tc.updatePayload)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			req := httptest.NewRequest(tc.method, tc.path, bytes.NewBuffer(reqBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Raw))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			finalHandler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Fatalf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			resp := new(data.User)
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if resp.UserName != newUserName {
				t.Errorf("Expected user name %s, got %s", newUserName, resp.UserName)
			}

		})

	}

}
