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

const (
	NameMinLen = 1
	NameMaxLen = 200
)
