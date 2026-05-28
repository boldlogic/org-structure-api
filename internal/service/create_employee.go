package service

import (
	"context"
	"fmt"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/boldlogic/packages/validate"
)

func (s *Service) CreateEmployee(ctx context.Context, emp models.Employee) (models.Employee, error) {
	trimmedName, err := validate.TrimAndLen(emp.FullName, models.NameMinLen, models.NameMaxLen)
	if err != nil {
		return models.Employee{}, fmt.Errorf("%w: %w: поле full_name", models.ErrValidation, err)
	}

	trimmedPosition, err := validate.TrimAndLen(emp.Position, models.NameMinLen, models.NameMaxLen)
	if err != nil {
		return models.Employee{}, fmt.Errorf("%w: %w: поле position", models.ErrValidation, err)
	}

	emp.FullName = trimmedName
	emp.Position = trimmedPosition

	return s.repo.CreateEmployee(ctx, emp)
}
