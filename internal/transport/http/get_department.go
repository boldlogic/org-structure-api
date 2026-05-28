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
			return nil, fmt.Sprintf("depth=%q, разрешено число от 1 до 5", rawDepth), models.ErrValidation
		}
	}
	includeFlag := true

	rawIncludeFlag := r.URL.Query().Get("include_employees")
	if rawIncludeFlag != "" {
		includeFlag, err = strconv.ParseBool(rawIncludeFlag)
		if err != nil {
			return nil, fmt.Sprintf("include_employees=%q, разрешено true или false", rawIncludeFlag), models.ErrValidation
		}
	}

	dep, children, emp, err := h.service.GetDepartment(r.Context(), id, depth, includeFlag)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, models.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}

	return departmentByIdToDTO(dep, children, emp), "", nil

}
