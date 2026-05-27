package server

import (
	"errors"
	"net/http"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/boldlogic/packages/transport/httputils"
)

func (h *Handler) createDepartment(r *http.Request) (any, string, error) {
	req, err := httputils.DecodeRequest[createDepartmentDTO](r)
	if err != nil {
		err = mapTransportErrors(err)
		return nil, err.Error(), err
	}
	dept, err := h.service.CreateDepartment(r.Context(), req.Name, req.ParentID)
	if err != nil {

		return nil, err.Error(), err
	}
	return departmentToDto(dept), "", nil
}

func mapTransportErrors(err error) error {
	if errors.Is(err, httputils.ErrUnsupportedMediaType) || errors.Is(err, httputils.ErrRequestEntityTooLarge) {
		return err
	}
	return models.ErrValidation
}
