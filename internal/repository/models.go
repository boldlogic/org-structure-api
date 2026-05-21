package repository

import (
	"time"
)

type department struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	ParentID  *int64    `gorm:"column:parent_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}
