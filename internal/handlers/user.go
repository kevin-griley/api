package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kevin-griley/api/internal/data"
	"github.com/kevin-griley/api/internal/middleware"
)

type PostUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary			Create a new user
// @Description		Create a new user
// @Tags			User
// @Accept			json
// @Produce			json
// @Param			body	body		PostUserRequest	true	"Create User Request"
// @Success         200		{object}	data.User	"User"
// @Failure         400		{object} 	ApiError	"Bad Request"
// @Router			/user	[post]
func HandlePostUser(w http.ResponseWriter, r *http.Request) *ApiError {
	ctx := r.Context()

	registerReq := new(PostUserRequest)
	if err := json.NewDecoder(r.Body).Decode(registerReq); err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	user, err := data.CreateRequest(registerReq.Email, registerReq.Password)
	if err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	resp, err := data.CreateUser(ctx, user)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// @Summary			Get user by apiKey
// @Description		Get user by apiKey
// @Tags			User
// @Security 		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			id			path	string	true	"User ID"
// @Success         200			{object}	data.User	"User"
// @Failure         400			{object} 	ApiError	"Bad Request"
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

type PatchUserRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// @Summary			Patch user by apiKey
// @Description		Patch user by apiKey
// @Tags			User
// @Security 		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			body		body		PatchUserRequest	true	"Patch User Request"
// @Success         200			{object}	data.User	"User"
// @Failure         400			{object} 	ApiError	"Bad Request"
// @Router			/user/me	[patch]
func HandlePatchUser(w http.ResponseWriter, r *http.Request) *ApiError {
	ctx := r.Context()

	patchReq := new(PatchUserRequest)
	if err := json.NewDecoder(r.Body).Decode(patchReq); err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	user, err := data.UpdateRequest(patchReq.UserName, patchReq.Password)
	if err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return &ApiError{http.StatusBadRequest, "Invalid user id"}
	}

	user.ID = userID

	resp, err := data.UpdateUser(ctx, user)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}

	return WriteJSON(w, http.StatusOK, resp)
}
