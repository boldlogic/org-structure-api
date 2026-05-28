package server

import (
	"errors"
	"net/http"

	"github.com/boldlogic/org-structure-api/internal/models"
)

func (h *Handler) deleteDepartment(r *http.Request) (any, string, error) {
	id, err := getDepartmentId(r)
	if err != nil {
		return nil, err.Error(), err
	}

	if err := h.service.DeleteDepartment(r.Context(), id); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}

	return nil, "", nil
}
