package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/boldlogic/packages/validate"
)

func (s *Service) CreateDepartment(ctx context.Context, name string, parentID *int64) (models.Department, error) {
	trimmed, err := validate.TrimAndLen(name, models.NameMinLen, models.NameMaxLen)
	if err != nil {
		return models.Department{}, fmt.Errorf("%w: %w: название департамента", models.ErrValidation, err)
	}

	out, err := s.repo.CreateDepartment(ctx, trimmed, parentID)
	if err != nil {
		if errors.Is(err, models.ErrConflict) {
			return models.Department{}, fmt.Errorf("%w: подразделение с названием %s уже существует у родителя", err, trimmed)
		}
		if errors.Is(err, models.ErrParentNotFound) {
			return models.Department{}, fmt.Errorf("%w: %w", models.ErrBusinessValidation, err)
		}
		return models.Department{}, err
	}
	return out, nil
}
