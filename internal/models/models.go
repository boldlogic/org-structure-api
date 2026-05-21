package models

import (
	"errors"
	"time"
)

type Department struct {
	ID        int64
	Name      string
	ParentID  *int64
	CreatedAt time.Time
}

var (
	ErrValidation = errors.New("некорректные входные данные")
)
