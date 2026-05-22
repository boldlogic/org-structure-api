package service

import (
	"context"
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
		return models.Department{}, err
	}
	return out, nil
}
