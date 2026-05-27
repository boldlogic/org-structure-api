package repository

import "context"

const (
	insert = `
	INSERT INTO org.employees(
	id, department_id, full_name, position, hired_at, created_at)
	VALUES (?, ?, ?, ?, ?, ?);`
)

func (r *Repo) insertEmployee(ctx context.Context, name string, parentID *int64) (result department, err error) {
	defer func() { r.logWrapper("insertEmployee", err) }()

	err = r.db.WithContext(ctx).Raw(insertNotExists, name, parentID).Scan(&result).Error
	return result, err
}
