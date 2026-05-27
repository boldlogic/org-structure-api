package repository

import (
	"context"
	"errors"

	"github.com/boldlogic/org-structure-api/internal/models"
	"gorm.io/gorm"
)

const (
	selectDepartment = `
		select
			id,
			name,
			parent_id,
			created_at
		FROM
			org.departments d
		where
			id = ?
		`
	selectDepartmentsByParentId = `
		WITH RECURSIVE tree AS (
			SELECT
				id,
				name,
				parent_id,
				created_at,
				1 AS level
			FROM
				org.departments
			WHERE
				parent_id = ?
			UNION ALL
			SELECT
				d.id,
				d.name,
				d.parent_id,
				d.created_at,
				tree.level + 1 AS level
			FROM
				org.departments d
				JOIN tree ON d.parent_id = tree.id
			WHERE
				tree.level < ?
		)
		SELECT
			id,
			name,
			parent_id,
			created_at
		FROM
			tree
		ORDER BY
			level,
			id
`
)

func (r *Repo) SelectDepartmentById(ctx context.Context, id int64) (result models.Department, err error) {
	defer func() { r.logWrapper("SelectDepartmentById", err) }()
	var row department
	err = r.db.WithContext(ctx).Raw(selectDepartment, id).Scan(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Department{}, models.ErrNotFound
		}
		return models.Department{}, err
	}
	if row.ID == 0 {
		return models.Department{}, models.ErrNotFound
	}

	return toDepartment(row), nil
}

func (r *Repo) SelectChildrenDepartments(ctx context.Context, parentId int64, depth int) (result []models.Department, err error) {
	defer func() { r.logWrapper("SelectChildrenDepartments", err) }()
	var rows []department
	err = r.db.WithContext(ctx).Raw(selectDepartmentsByParentId, parentId, depth).Scan(&rows).Error

	if err != nil {
		return nil, err
	}
	result = toDepartments(rows)
	return result, nil
}
