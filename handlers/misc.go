package handlers

import (
	"fmt"
	"net/http"
)

func GetID(r *http.Request) (string, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return idStr, fmt.Errorf("id is required")
	}
	return idStr, nil
}
