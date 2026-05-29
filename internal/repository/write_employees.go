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
const updateEmployeeDepartmentID = `
update org.employees
set department_id=?
where department_id=?`

func (r *Repo) ReassignAndDelete(ctx context.Context, departmentId int64, newDepartmentId int64) (err error) {
	defer func() { r.logWrapper("ReassignAndDelete", err) }()

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = tx.Exec(updateEmployeeDepartmentID, newDepartmentId, departmentId).Error
	if err != nil {
		return err
	}
	err = tx.Exec(updateDepartmentParentId, newDepartmentId, departmentId).Error
	if err != nil {
		return err
	}
	res := tx.Exec(deleteDepartment, departmentId)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return models.ErrNotFound
	}
	err = tx.Commit().Error
	return err
}

func (r *Repo) CreateEmployee(ctx context.Context, emp models.Employee) (result models.Employee, err error) {
	defer func() { r.logWrapper("CreateEmployee", err) }()

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
