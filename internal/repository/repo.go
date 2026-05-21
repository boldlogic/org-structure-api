package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewRepo(db *gorm.DB, logger *zap.Logger) *Repo {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Repo{
		db:     db,
		logger: logger,
	}
}
