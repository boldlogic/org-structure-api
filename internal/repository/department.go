package repository

import (
	"context"
	"errors"

	"github.com/boldlogic/org-structure-api/internal/models"
	"gorm.io/gorm"
)

const (
	insertNotExists = `
		WITH src AS (SELECT ?::varchar(200) AS name, ?::integer AS parent_id)
		INSERT INTO departments (name, parent_id)
		SELECT name, parent_id FROM src
		WHERE NOT EXISTS (
			SELECT 1 FROM departments d
			WHERE d.name = src.name
			AND (
				(d.parent_id IS NULL AND src.parent_id IS NULL)
				OR d.parent_id = src.parent_id
			)
		)
		RETURNING id, name, parent_id, created_at
`
)

func (r *Repo) CreateDepartment(ctx context.Context, name string, parentID *int64) (result models.Department, err error) {
	defer func() { r.logWrapper("CreateDepartment", err) }()

	var row department
	err = r.db.WithContext(ctx).Raw(insertNotExists, name, parentID).Scan(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Department{}, models.ErrConflict
		}
		return models.Department{}, err
	}

	if row.ID == 0 {
		return models.Department{}, models.ErrConflict
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
