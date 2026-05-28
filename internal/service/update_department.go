package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/boldlogic/packages/validate"
)

func (s *Service) UpdateDepartment(ctx context.Context, id int64, name *string, parentID *int64, parentIDSet bool) (models.Department, error) {
	if name == nil && !parentIDSet {
		return models.Department{}, fmt.Errorf("%w: не передано ни одно поле для обновления", models.ErrValidation)
	}

	if name != nil {
		trimmed, err := validate.TrimAndLen(*name, models.NameMinLen, models.NameMaxLen)
		if err != nil {
			return models.Department{}, fmt.Errorf("%w: %w: название департамента", models.ErrValidation, err)
		}
		name = &trimmed
	}

	if parentIDSet && parentID != nil && *parentID == id {
		return models.Department{}, fmt.Errorf("%w: parent_id не может быть равен id", models.ErrValidation)
	}

	out, err := s.repo.UpdateDepartment(ctx, id, name, parentID, parentIDSet)
	if err != nil {
		if errors.Is(err, models.ErrCycle) {
			return models.Department{}, fmt.Errorf(
				"%w: нельзя сделать подразделение %d потомком своего поддерева",
				models.ErrConflict, id,
			)
		}
		if errors.Is(err, models.ErrConflict) {
			return models.Department{}, fmt.Errorf("%w: подразделение с таким названием уже существует у родителя", err)
		}
		if errors.Is(err, models.ErrParentNotFound) && parentID != nil {
			return models.Department{}, fmt.Errorf("%w: %w %d", models.ErrBusinessValidation, err, *parentID)
		}
		return models.Department{}, err
	}
	return out, nil
}
