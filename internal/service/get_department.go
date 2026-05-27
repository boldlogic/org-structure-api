package service

import (
	"context"
	"fmt"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/boldlogic/packages/validate"
)

func (s *Service) GetDepartment(ctx context.Context, id int64, depth int) (models.Department, []models.Department, error) {

	err := validate.CheckIntInRange(depth, 1, 5)
	if err != nil {
		return models.Department{}, nil, fmt.Errorf("%w: некорректное значение depth=%d: %w", models.ErrBusinessValidation, depth, err)
	}

	dep, err := s.repo.SelectDepartmentById(ctx, id)
	if err != nil {
		return models.Department{}, nil, err
	}
	var children []models.Department

	children, err = s.repo.SelectChildrenDepartments(ctx, dep.ID, depth)
	if err != nil {
		return models.Department{}, nil, err
	}

	return dep, children, nil
}
