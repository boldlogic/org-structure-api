package service

import (
	"context"

	"github.com/boldlogic/org-structure-api/internal/models"
)

type repository interface {
	CreateDepartment(ctx context.Context, name string, parentID *int64) (models.Department, error)
	CreateEmployee(ctx context.Context, emp models.Employee) (models.Employee, error)
	SelectDepartmentById(ctx context.Context, id int64) (result models.Department, err error)
	SelectChildrenDepartments(ctx context.Context, parentId int64, depth int) ([]models.Department, error)
	UpdateDepartment(ctx context.Context, id int64, name *string, parentID *int64, parentIDSet bool) (models.Department, error)
	DeleteDepartment(ctx context.Context, id int64) error
	SelectEmployeesDepartmentById(ctx context.Context, id int64) (result []models.Employee, err error)
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{repo: repo}
}
