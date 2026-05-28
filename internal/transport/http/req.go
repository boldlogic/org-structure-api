package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/boldlogic/packages/validate"
)

func getDepartmentId(r *http.Request) (int64, error) {
	raw := r.PathValue("id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("%w: id=%s", models.ErrValidation, raw)
	}
	err = validate.CheckIntInRange(int(id), 1, 2147483647)
	if err != nil {
		return 0, fmt.Errorf("%w: %w id=%s", models.ErrValidation, err, raw)
	}
	return id, nil

}
