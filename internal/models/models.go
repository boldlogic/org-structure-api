package models

import (
	"time"
)

type Department struct {
	ID        int64
	Name      string
	ParentID  *int64
	CreatedAt time.Time
}

type Employee struct {
	ID           int64
	DepartmentID int64
	FullName     string
	Position     string
	HiredAt      *time.Time
	CreatedAt    time.Time
}

const (
	NameMinLen = 1
	NameMaxLen = 200
)
