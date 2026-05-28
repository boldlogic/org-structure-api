package service

import (
	"context"
	"testing"
	"time"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testRepo struct {
	mock.Mock
}

func (r *testRepo) CreateDepartment(ctx context.Context, name string, parentID *int64) (models.Department, error) {
	args := r.Called(ctx, name, parentID)

	dep, _ := args.Get(0).(models.Department)
	return dep, args.Error(1)
}

func (r *testRepo) SelectDepartmentById(_ context.Context, id int64) (models.Department, error) {
	return models.Department{}, nil
}

func (r *testRepo) SelectChildrenDepartments(context.Context, int64, int) ([]models.Department, error) {
	return nil, nil
}

func (r *testRepo) UpdateDepartment(context.Context, int64, *string, *int64, bool) (models.Department, error) {
	return models.Department{}, nil
}

func (r *testRepo) CreateEmployee(_ context.Context, _ models.Employee) (models.Employee, error) {
	return models.Employee{}, nil
}

func (r *testRepo) DeleteDepartment(_ context.Context, _ int64) error {
	return nil
}

func (r *testRepo) SelectEmployeesDepartmentById(ctx context.Context, id int64) (result []models.Employee, err error) {
	return nil, nil
}

func Test_CreateDepartment(t *testing.T) {
	parentID := int64(1)
	otherParentID := int64(5)
	сreated := time.Date(2026, 05, 27, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		inName     string
		inParent   *int64
		repoName   string
		repoParent *int64
		repoDep    models.Department
		repoErr    error
		wantDep    models.Department
		wantErr    error
	}{
		{
			name:     "успешное_создание",
			inName:   "новый департамент",
			repoName: "новый департамент",
			repoDep: models.Department{
				ID:        4,
				Name:      "новый департамент",
				CreatedAt: сreated,
			},
			wantDep: models.Department{
				ID:        4,
				Name:      "новый департамент",
				CreatedAt: сreated,
			},
		},
		{
			name:       "успешное_создание_с_parent_id",
			inName:     "новый департамент 2",
			inParent:   &parentID,
			repoName:   "новый департамент 2",
			repoParent: &parentID,
			repoDep: models.Department{
				ID:        4,
				Name:      "новый департамент 2",
				ParentID:  &parentID,
				CreatedAt: сreated,
			},
			wantDep: models.Department{
				ID:        4,
				Name:      "новый департамент 2",
				ParentID:  &parentID,
				CreatedAt: сreated,
			},
		},
		{
			name:    "пустое_имя",
			inName:  "",
			wantErr: models.ErrValidation,
		},
		{
			name:    "пустое_имя_с_пробелом",
			inName:  " ",
			wantErr: models.ErrValidation,
		},
		{
			name:     "обрезает_пробелы",
			inName:   " новый департамент ",
			repoName: "новый департамент",
			repoDep: models.Department{
				ID:        4,
				Name:      "новый департамент",
				CreatedAt: сreated,
			},
			wantDep: models.Department{
				ID:        4,
				Name:      "новый департамент",
				CreatedAt: сreated,
			},
		},
		{
			name:     "конфликт",
			inName:   " Подразделение ",
			repoName: "Подразделение",
			repoErr:  models.ErrConflict,
			wantErr:  models.ErrConflict,
		},
		{
			name:       "parent_id_не_существует",
			inName:     "Подразделение",
			inParent:   &otherParentID,
			repoName:   "Подразделение",
			repoParent: &otherParentID,
			repoErr:    models.ErrParentNotFound,
			wantErr:    models.ErrBusinessValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(testRepo)
			if tt.repoName != "" {
				repo.On("CreateDepartment", mock.Anything, tt.repoName, tt.repoParent).
					Return(tt.repoDep, tt.repoErr).
					Once()
			}

			svc := NewService(repo)
			got, err := svc.CreateDepartment(context.Background(), tt.inName, tt.inParent)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, models.Department{}, got)
				if tt.repoName == "" {
					repo.AssertNotCalled(t, "CreateDepartment")
				}
				repo.AssertExpectations(t)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantDep, got)
			repo.AssertExpectations(t)
		})
	}
}
