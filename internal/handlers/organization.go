package handlers

import (
	"net/http"

	"github.com/kevin-griley/api/internal/data"
)

type PostOrganizationRequest struct {
	Name             string                `json:"name"`
	Address          string                `json:"address"`
	ContactInfo      string                `json:"contact_info"`
	OrganizationType data.OrganizationType `json:"organization_type"`
}

// @Summary			Create a new organization
// @Description		Create a new organization
// @Tags			Organization
// @Accept			json
// @Produce			json
// @Param			body	body		PostOrganizationRequest	true	"Create Organization Request"
// @Success         200		{object}	data.Organization	"Organization"
// @Failure         400		{object} 	ApiError	"Bad Request"
// @Router			/organization	[post]
func HandlePostOrganization(w http.ResponseWriter, r *http.Request) *ApiError {
	ctx := r.Context()

	store, ok := data.GetStore(ctx)
	if !ok {
		return &ApiError{http.StatusInternalServerError, "no database store in context"}
	}

	postReq := new(PostOrganizationRequest)
	if err := DecodeJSONRequest(r, postReq, 1<<20); err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	org, err := store.Organization.CreateRequest(postReq.Name, postReq.Address, postReq.ContactInfo, postReq.OrganizationType)
	if err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	resp, err := store.Organization.CreateOrganization(org)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// @Summary			Get organization by ID
// @Description		Get organization by ID
// @Tags			Organization
// @Security 		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			id	path	string	true	"Organization ID"
// @Success         200			{object}	data.Organization	"Organization"
// @Failure         400			{object} 	ApiError	"Bad Request"
// @Router			/organization/{id}	[get]
func HandleGetOrganizationByID(w http.ResponseWriter, r *http.Request) *ApiError {
	ctx := r.Context()

	store, ok := data.GetStore(ctx)
	if !ok {
		return &ApiError{http.StatusInternalServerError, "no database store in context"}
	}

	orgId, err := GetPathID(r)
	if err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	org, err := store.Organization.GetOrganizationByID(orgId)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}

	return WriteJSON(w, http.StatusOK, org)
}

type PatchOrganizationRequest struct {
	name             string
	uniqueURL        string
	address          string
	contactInfo      string
	organizationType data.OrganizationType
}

// @Summary			Patch organization by ID
// @Description		Patch organization by ID
// @Tags			Organization
// @Security 		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			id	path	string	true	"Organization ID"
// @Param			body	body		PatchOrganizationRequest	true	"Patch Organization Request"
// @Success         200			{object}	data.Organization	"Organization"
// @Failure         400			{object} 	ApiError	"Bad Request"
// @Router			/organization/{id}	[patch]
func HandlePatchOrganizationByID(w http.ResponseWriter, r *http.Request) *ApiError {
	ctx := r.Context()

	store, ok := data.GetStore(ctx)
	if !ok {
		return &ApiError{http.StatusInternalServerError, "no database store in context"}
	}

	patchReq := new(PatchOrganizationRequest)
	if err := DecodeJSONRequest(r, patchReq, 1<<20); err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	org, err := store.Organization.UpdateRequest(
		patchReq.name,
		patchReq.uniqueURL,
		patchReq.address,
		patchReq.contactInfo,
		patchReq.organizationType,
	)
	if err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	orgId, err := GetPathID(r)
	if err != nil {
		return &ApiError{http.StatusBadRequest, err.Error()}
	}

	org.ID = orgId

	resp, err := store.Organization.UpdateOrganization(org)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}

	return WriteJSON(w, http.StatusOK, resp)

}
