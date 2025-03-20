package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kevin-griley/api/internal/data"
)

type PostAuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PostAuthResponse struct {
	Token string `json:"token"`
}

// @Summary			Retrive token for bearer authentication
// @Description		Retrive token for bearer authentication
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			body	body		PostAuthRequest	true	"Login Request"
// @Success			200		{object}	PostAuthResponse	"Token Response"
// @Failure			400		{object} 	ApiError	"Bad Request"
// @Failure			401		{object} 	ApiError	"Unauthorized"
// @Router			/login	[post]
func HandlePostLogin(w http.ResponseWriter, r *http.Request) *ApiError {
	ctx := r.Context()
	rBody := new(PostAuthRequest)
	if err := json.NewDecoder(r.Body).Decode(rBody); err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	user, err := data.GetUserByEmail(ctx, rBody.Email)
	if err != nil {
		return &ApiError{http.StatusUnauthorized, "invalid user or password"}
	}

	if user.IsDeleted {
		return &ApiError{http.StatusUnauthorized, "invalid user or password"}
	}

	if user.FailedLoginAttempts >= 10 && time.Since(user.UpdatedAt).Minutes() < 30 {
		return &ApiError{http.StatusUnauthorized, "account locked due to too many failed login attempts please try again later"}
	}

	if !user.ValidPassword(rBody.Password) {
		user.FailedLoginAttempts++

		_, err := data.UpdateUser(ctx, user)
		if err != nil {
			return &ApiError{http.StatusInternalServerError, err.Error()}
		}

		return &ApiError{http.StatusUnauthorized, "invalid user or password"}
	}

	user.FailedLoginAttempts = 1
	user.LastLogin = time.Now().UTC()
	user, err = data.UpdateUser(ctx, user)

	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}

	tokenString, err := CreateJWT(user)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}

	return WriteJSON(w, http.StatusOK, PostAuthResponse{Token: tokenString})
}

func CreateJWT(user *data.User) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "mycartage",
		Subject:   user.ID.String(),
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))

}
