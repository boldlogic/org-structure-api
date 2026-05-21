package service

import (
	"context"

	"github.com/boldlogic/org-structure-api/internal/models"
)

type repository interface {
	CreateDepartment(ctx context.Context, name string, parentID *int64) (models.Department, error)
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateDepartment(ctx context.Context, name string, parentID *int64) (models.Department, error) {
	return s.repo.CreateDepartment(ctx, name, parentID)
}
