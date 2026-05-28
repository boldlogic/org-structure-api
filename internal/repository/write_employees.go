package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

const insertEmployee = `
	INSERT INTO org.employees (department_id, full_name, position, hired_at)
	VALUES (?, ?, ?, ?)
	RETURNING id, department_id, full_name, position, hired_at, created_at
`

func (r *Repo) CreateEmployee(ctx context.Context, emp models.Employee) (result models.Employee, err error) {
	defer func() { r.logWrapper("CreateEmployee", err) }()

	_, err = r.SelectDepartmentById(ctx, emp.DepartmentID)
	if err != nil {
		return models.Employee{}, err
	}

	var row employee
	err = r.db.WithContext(ctx).Raw(insertEmployee, emp.DepartmentID, emp.FullName, emp.Position, emp.HiredAt).Scan(&row).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return models.Employee{}, models.ErrNotFound
		}
		return models.Employee{}, err
	}

	if row.ID == 0 {
		return models.Employee{}, fmt.Errorf("запись не была создана")
	}

	return toEmployee(row), nil
}
