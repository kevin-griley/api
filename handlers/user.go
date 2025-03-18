package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kevin-griley/api/data"
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
	registerReq := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(registerReq); err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	user, err := data.NewUser(registerReq.Email, registerReq.Password)
	if err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	userID, err := data.CreateUser(user)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}

	user.ID = userID
	return WriteJSON(w, http.StatusOK, user)

}

// @Summary			Get User by ID
// @Description		Get a user by ID
// @Tags			User
// @Security 		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			id		path	string	true	"User ID"
// @Success         200		{object}	data.User	"User"
// @Failure         400		{object} 	ApiError	"Bad Request"
// @Router			/user/{id}	[get]
func HandleGetUser(w http.ResponseWriter, r *http.Request) *ApiError {
	id, err := GetID(r)
	if err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	user, err := data.GetUserByID(id)
	if err != nil {
		return &ApiError{http.StatusNotFound, err.Error()}
	}

	return WriteJSON(w, http.StatusOK, user)
}
