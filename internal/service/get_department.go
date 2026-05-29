package service

import (
	"context"
	"fmt"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/boldlogic/packages/validate"
	"golang.org/x/sync/errgroup"
)

func (s *Service) GetDepartment(ctx context.Context, id int64, depth int, includeEmployeeFlag bool) (models.Department, []models.Department, []models.Employee, error) {

	err := validate.CheckIntInRange(depth, 1, 5)
	if err != nil {
		return models.Department{}, nil, nil, fmt.Errorf("%w: depth=%d: %w", models.ErrBusinessValidation, depth, err)
	}

	dep, err := s.repo.SelectDepartmentById(ctx, id)
	if err != nil {
		return models.Department{}, nil, nil, err
	}

	var emp []models.Employee
	var children []models.Department

	g, gCtx := errgroup.WithContext(ctx)
	if includeEmployeeFlag {
		g.Go(func() error {
			var err error
			emp, err = s.repo.SelectEmployeesDepartmentById(gCtx, id)
			return err

		})

	}
	g.Go(func() error {
		var err error
		children, err = s.repo.SelectChildrenDepartments(gCtx, id, depth)
		return err

	})

	err = g.Wait()

	if err != nil {
		return models.Department{}, nil, nil, err
	}

	return dep, children, emp, nil
}
