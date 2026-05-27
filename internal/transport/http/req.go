package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/boldlogic/org-structure-api/internal/models"
)

func getDepartmentId(r *http.Request) (int64, error) {
	raw := r.PathValue("id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("%w: некорректное значение id: %s", models.ErrValidation, raw)
	}
	return id, nil

}
