package server

import (
	"errors"
	"net/http"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/boldlogic/packages/transport/httputils"
)

type departmentReqDTO struct {
	Name     string `json:"name" validate:"required,min=1,max=200"`
	ParentID *int64 `json:"parent_id" validate:"omitempty"`
}

func (h *Handler) createDepartment(r *http.Request) (any, string, error) {
	req, err := httputils.DecodeRequest[departmentReqDTO](r)
	if err != nil {
		if errors.Is(err, httputils.ErrUnsupportedMediaType) || errors.Is(err, httputils.ErrRequestEntityTooLarge) {
			return nil, err.Error(), err
		}
		return nil, err.Error(), models.ErrValidation
	}
	dept, err := h.service.CreateDepartment(r.Context(), req.Name, req.ParentID)
	if err != nil {
		return nil, "", err
	}
	return dept, "", nil
}
