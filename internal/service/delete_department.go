package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/boldlogic/org-structure-api/internal/models"
)

func (s *Service) DeleteDepartment(ctx context.Context, id int64, mode string, reassignToDepartment *int64) error {
	var err error

	switch mode {
	case "reassign":
		if reassignToDepartment == nil {
			return fmt.Errorf("%w: новый департамент обязателен при reassign", models.ErrBusinessValidation)
		}

		if *reassignToDepartment == id {
			return fmt.Errorf("%w: новый департамент не может быть равен id удаляемого департамента", models.ErrBusinessValidation)
		}
		_, err = s.repo.SelectDepartmentById(ctx, *reassignToDepartment)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return fmt.Errorf("%w: департамента с id=%d не существует", models.ErrBusinessValidation, *reassignToDepartment)

			}
			return err
		}
		children, err := s.repo.SelectChildrenDepartments(ctx, id, 1)
		if err != nil {
			return err
		}
		for _, c := range children {
			if c.ID == *reassignToDepartment {
				return fmt.Errorf("%w: дочерний департамент %d не может быть новым департаментом", models.ErrBusinessValidation, c.ID)
			}

		}

		err = s.repo.ReassignAndDelete(ctx, id, *reassignToDepartment)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return fmt.Errorf("%w: департамента с id=%d не существует", models.ErrBusinessValidation, *reassignToDepartment)

			}
			return err
		}

	case "cascade":
		err = s.repo.DeleteDepartment(ctx, id)
	default:
		return fmt.Errorf("%w: mode=%q, разрешено cascade или reassign", models.ErrBusinessValidation, mode)
	}

	return err
}
