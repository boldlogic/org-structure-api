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

const (
	NameMinLen = 1
	NameMaxLen = 200
)

var (
	ErrValidation = errors.New("некорректные входные данные")
	ErrConflict   = errors.New("запись с таким ключом уже существует")
)
