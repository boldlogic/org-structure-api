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
		INSERT INTO
			org.departments (name, parent_id)
		SELECT
			name,
			parent_id
		FROM
			src
		WHERE
			NOT EXISTS (
				SELECT
					1
				FROM
					org.departments d
				WHERE
					d.name = src.name
					AND (
						(
							d.parent_id IS NULL
							AND src.parent_id IS NULL
						)
						OR d.parent_id = src.parent_id
					)
			)
			AND (
				src.parent_id IS NULL
				or exists(
					select
						1
					from
						org.departments e
					where
						e.id = src.parent_id
				)
			)
		RETURNING id, name, parent_id, created_at
`
	updateDepartment = `	
		UPDATE
			org.departments d
		SET
			name = COALESCE(?, d.name),
			parent_id = CASE WHEN ? THEN ? ELSE d.parent_id END
		WHERE
			d.id = ? RETURNING d.id,
			d.name,
			d.parent_id,
			d.created_at		
`
	checkReason = `
		WITH q AS (
			SELECT
				? :: varchar(200) AS name,
				? :: integer AS parent_id
		),
		src as (
			SELECT
				name,
				parent_id,
				case when exists (
					select
						1
					from
						org.departments d
					where
						d.name = q.name
						AND (
							(
								d.parent_id IS NULL
								AND q.parent_id IS NULL
							)
							OR d.parent_id = q.parent_id
						)
				) then 'double_name' 
				when not exists (
					select
						1
					from
						org.departments e
					where
						e.id = q.parent_id
				) then 'not_existing_parent_id' 
				else 'ok' end as reason
			FROM
				q
		)
		select
			reason
		from
			src `
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

func toDepartment(r department) models.Department {
	return models.Department{
		ID:        r.ID,
		Name:      r.Name,
		ParentID:  r.ParentID,
		CreatedAt: r.CreatedAt,
	}
}

func toDepartments(r []department) []models.Department {
	var out = make([]models.Department, 0, len(r))

	for _, d := range r {
		out = append(out, toDepartment(d))
	}

	return out
}
