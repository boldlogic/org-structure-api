package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/boldlogic/org-structure-api/internal/models"
)

func (h *Handler) getDepartment(r *http.Request) (any, string, error) {
	id, err := getDepartmentId(r)
	if err != nil {
		return nil, err.Error(), err
	}
	rawDepth := r.URL.Query().Get("depth")
	depth := 1

	if rawDepth != "" {
		depth, err = strconv.Atoi(rawDepth)
		if err != nil {
			return nil, fmt.Sprintf("некорректное значение depth: %s", rawDepth), models.ErrValidation
		}
	}

	dep, children, err := h.service.GetDepartment(r.Context(), id, depth)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, models.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}

	return departmentByIdToDTO(dep, children), "", nil

}
