package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kevin-griley/api/internal/data"
	"github.com/kevin-griley/api/internal/middleware"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary			Create User
// @Description		Create a new user
// @Tags			User
// @Accept			json
// @Produce			json
// @Param			body	body		CreateUserRequest	true	"Create User Request"
// @Success         200		{object}	data.User	"User"
// @Failure         400		{object} 	ApiError	"Bad Request"
// @Router			/user	[post]
func HandlePostUser(w http.ResponseWriter, r *http.Request) *ApiError {
	ctx := r.Context()

	registerReq := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(registerReq); err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	user, err := data.NewUser(registerReq.Email, registerReq.Password)
	if err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	userID, err := data.CreateUser(ctx, user)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}

	user.ID = userID
	return WriteJSON(w, http.StatusOK, user)

}

// @Summary			Get User by apiKey
// @Description		Get a user by apiKey
// @Tags			User
// @Security 		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			id		path	string	true	"User ID"
// @Success         200		{object}	data.User	"User"
// @Failure         400		{object} 	ApiError	"Bad Request"
// @Router			/user/me	[get]
func HandleGetUserByKey(w http.ResponseWriter, r *http.Request) *ApiError {
	ctx := r.Context()

	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return &ApiError{http.StatusBadRequest, "Invalid user id"}
	}

	user, err := data.GetUserByID(ctx, userID)
	if err != nil {
		return &ApiError{http.StatusNotFound, err.Error()}
	}

	return WriteJSON(w, http.StatusOK, user)
}
