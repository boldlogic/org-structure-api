package server

import (
	"errors"
	"net/http"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/boldlogic/packages/transport/httputils"
)

func (h *Handler) updateDepartment(r *http.Request) (any, string, error) {
	id, err := getDepartmentId(r)
	if err != nil {
		return nil, err.Error(), err
	}

	req, err := httputils.DecodeRequest[updateDepartmentDTO](r)
	if err != nil {
		if errors.Is(err, httputils.ErrUnsupportedMediaType) || errors.Is(err, httputils.ErrRequestEntityTooLarge) {
			return nil, err.Error(), err
		}
		return nil, err.Error(), models.ErrValidation
	}

	if req.Name == nil {
		return nil, "name не может быть null", models.ErrValidation
	}

	var name *string
	if *req.Name != nil {
		name = *req.Name
	}

	parentIDSet := req.ParentID == nil || *req.ParentID != nil
	var parentID *int64
	if req.ParentID != nil && *req.ParentID != nil {
		parentID = *req.ParentID
	}

	dept, err := h.service.UpdateDepartment(r.Context(), id, name, parentID, parentIDSet)
	if err != nil {
		if errors.Is(err, models.ErrValidation) || errors.Is(err, models.ErrConflict) || errors.Is(err, models.ErrNotFound) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}
	return departmentToDto(dept), "", nil
}
