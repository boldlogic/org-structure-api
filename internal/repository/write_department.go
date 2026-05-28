package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	insertNotExists = `
		WITH src AS (
			SELECT
				?::varchar(200) AS name,
				?::integer AS parent_id
		)
		INSERT INTO org.departments (name, parent_id)
		SELECT
			name,
			parent_id
		FROM src
		WHERE
			NOT EXISTS (
				SELECT 1
				FROM org.departments d
				WHERE
					d.name = src.name
					AND (
						(d.parent_id IS NULL AND src.parent_id IS NULL)
						OR d.parent_id = src.parent_id
					)
			)
			AND (
				src.parent_id IS NULL
				OR EXISTS (
					SELECT 1
					FROM org.departments e
					WHERE e.id = src.parent_id
				)
			)
		RETURNING
			id,
			name,
			parent_id,
			created_at
	`
	updateDepartment = `
		UPDATE org.departments d
		SET
			name = COALESCE(?, d.name),
			parent_id = CASE
				WHEN ? THEN ?
				ELSE d.parent_id
			END
		WHERE
			d.id = ?
		RETURNING
			d.id,
			d.name,
			d.parent_id,
			d.created_at
	`

	updateDepartmentParentId = `
		UPDATE org.departments d
		SET 			parent_id = ?
		WHERE
			d.parent_id  = ?
	`
	checkReason = `
		WITH q AS (
			SELECT
				?::varchar(200) AS name,
				?::integer AS parent_id
		),
		src AS (
			SELECT
				name,
				parent_id,
				CASE
					WHEN EXISTS (
						SELECT 1
						FROM org.departments d
						WHERE
							d.name = q.name
							AND (
								(d.parent_id IS NULL AND q.parent_id IS NULL)
								OR d.parent_id = q.parent_id
							)
					) THEN 'double_name'
					WHEN NOT EXISTS (
						SELECT 1
						FROM org.departments e
						WHERE e.id = q.parent_id
					) THEN 'not_existing_parent_id'
					ELSE 'ok'
				END AS reason
			FROM q
		)
		SELECT reason
		FROM src
	`
	checkDepartmentCycle = `
		WITH RECURSIVE sub AS (
			SELECT
				id
			FROM
				org.departments
			WHERE
				id = ?::integer
			UNION ALL
			SELECT
				d.id
			FROM
				org.departments d
				INNER JOIN sub ON d.parent_id = sub.id
		)
		SELECT
			CASE
				WHEN EXISTS (
					SELECT 1
					FROM sub
					WHERE id = ?
				) THEN 'cycle'
				ELSE 'ok'
			END AS reason
	`
)

func (r *Repo) insertDepartment(ctx context.Context, name string, parentID *int64) (result department, err error) {
	defer func() { r.logWrapper("insertDepartment", err) }()

	err = r.db.WithContext(ctx).Raw(insertNotExists, name, parentID).Scan(&result).Error
	return result, err
}

func (r *Repo) CreateDepartment(ctx context.Context, name string, parentID *int64) (result models.Department, err error) {
	defer func() { r.logWrapper("CreateDepartment", err) }()

	var row department
	row, err = r.insertDepartment(ctx, name, parentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows) {
			goto NoRows
		}

		return models.Department{}, err
	}

NoRows:
	if row.ID == 0 {
		var reason string
		err = r.db.WithContext(ctx).Raw(checkReason, name, parentID).Scan(&reason).Error
		var pgErr *pgconn.PgError
		if err != nil {
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return models.Department{}, models.ErrConflict
			}
			return models.Department{}, err
		}

		switch reason {
		case "double_name":
			return models.Department{}, models.ErrConflict
		case "not_existing_parent_id":
			return models.Department{}, models.ErrParentNotFound
		case "ok":
			row, err = r.insertDepartment(ctx, name, parentID)
		}

		if err != nil {
			return models.Department{}, err
		}
	}
	return toDepartment(row), nil
}

func (r *Repo) UpdateDepartment(ctx context.Context, id int64, name *string, parentID *int64, parentIDSet bool) (result models.Department, err error) {
	defer func() { r.logWrapper("UpdateDepartment", err) }()
	if name == nil && !parentIDSet {
		return models.Department{}, models.ErrValidation
	}

	if parentIDSet && parentID != nil {
		var reason string
		err = r.db.WithContext(ctx).Raw(checkDepartmentCycle, id, *parentID).Scan(&reason).Error
		if err != nil {
			return models.Department{}, err
		}
		if reason == "cycle" {
			return models.Department{}, models.ErrCycle
		}
	}

	var row department
	var pgErr *pgconn.PgError
	err = r.db.WithContext(ctx).Raw(updateDepartment, name, parentIDSet, parentID, id).Scan(&row).Error
	if err != nil {
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return models.Department{}, models.ErrConflict
			case "23503":
				return models.Department{}, models.ErrParentNotFound
			}
		}
		return models.Department{}, err
	}

	if row.ID == 0 {
		return models.Department{}, models.ErrNotFound
	}
	return toDepartment(row), nil
}

const deleteDepartment = `DELETE FROM org.departments WHERE id = ?`

func (r *Repo) DeleteDepartment(ctx context.Context, id int64) (err error) {
	defer func() { r.logWrapper("DeleteDepartment", err) }()

	res := r.db.WithContext(ctx).Exec(deleteDepartment, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}
