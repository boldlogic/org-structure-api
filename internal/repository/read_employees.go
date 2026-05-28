package repository

import (
	"context"

	"github.com/boldlogic/org-structure-api/internal/models"
)

const (
	selectEmployees = `
		SELECT id, department_id, full_name, position, hired_at, created_at
		FROM org.employees
		where
			department_id = ?
		`
)

func (r *Repo) SelectEmployeesDepartmentById(ctx context.Context, id int64) (result []models.Employee, err error) {
	defer func() { r.logWrapper("SelectEmployeesDepartmentById", err) }()
	var rows []employee
	err = r.db.WithContext(ctx).Raw(selectEmployees, id).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result = toEmployees(rows)
	return result, nil
}
