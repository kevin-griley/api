package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func GetID(r *http.Request) (uuid.UUID, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("id is required")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id: %s", idStr)
	}

	return id, nil
}
