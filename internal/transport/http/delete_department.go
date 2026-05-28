package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/boldlogic/packages/validate"
)

func (h *Handler) deleteDepartment(r *http.Request) (any, string, error) {
	id, err := getDepartmentId(r)
	if err != nil {
		return nil, err.Error(), err
	}

	rawMode := r.URL.Query().Get("mode")
	mode := strings.TrimSpace(rawMode)
	if mode == "" {
		return nil, "mode обязателен", models.ErrValidation
	}
	rawReassign := r.URL.Query().Get("reassign_to_department_id")
	var reassignTo *int64
	if strings.TrimSpace(rawReassign) != "" {
		parsed, err := strconv.ParseInt(rawReassign, 10, 64)
		if err != nil {
			return nil, fmt.Sprintf("reassign_to_department_id=%q, ожидается число", rawReassign), models.ErrValidation
		}
		err = validate.CheckIntInRange(int(parsed), 1, 2147483647)
		if err != nil {
			return nil, fmt.Sprintf("%v: reassign_to_department_id=%d", err, parsed), models.ErrValidation
		}
		reassignTo = &parsed
	}

	if err := h.service.DeleteDepartment(r.Context(), id, mode, reassignTo); err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, models.ErrBusinessValidation) {
			return nil, err.Error(), err
		}
		return nil, "", err
	}

	return nil, "", nil
}
