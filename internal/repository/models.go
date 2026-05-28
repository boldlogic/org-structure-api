package repository

import (
	"time"

	"github.com/boldlogic/org-structure-api/internal/models"
)

type department struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	ParentID  *int64    `gorm:"column:parent_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type employee struct {
	ID           int64      `gorm:"column:id;primaryKey"`
	DepartmentID int64      `gorm:"column:department_id"`
	FullName     string     `gorm:"column:full_name"`
	Position     string     `gorm:"column:position"`
	HiredAt      *time.Time `gorm:"column:hired_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
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

func toEmployee(r employee) models.Employee {
	return models.Employee{
		ID:           r.ID,
		DepartmentID: r.DepartmentID,
		FullName:     r.FullName,
		Position:     r.Position,
		HiredAt:      r.HiredAt,
		CreatedAt:    r.CreatedAt,
	}
}

func toEmployees(r []employee) []models.Employee {
	var out = make([]models.Employee, 0, len(r))
	for _, d := range r {
		out = append(out, toEmployee(d))
	}

	return out
}
