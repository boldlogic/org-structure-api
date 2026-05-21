package repository

import (
	"context"

	"github.com/boldlogic/org-structure-api/internal/models"
)

func (r *Repo) CreateDepartment(ctx context.Context, name string, parentID *int64) (result models.Department, err error) {
	defer func() { r.logWrapper("CreateDepartment", err) }()

	row := department{Name: name, ParentID: parentID}
	if err = r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return models.Department{}, err
	}
	result = toDepartment(row)
	return result, nil
}

func toDepartment(r department) models.Department {
	return models.Department{
		ID:        r.ID,
		Name:      r.Name,
		ParentID:  r.ParentID,
		CreatedAt: r.CreatedAt,
	}
}
