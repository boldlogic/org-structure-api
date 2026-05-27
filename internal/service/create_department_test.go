package service

import (
	"context"
	"testing"
	"time"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/stretchr/testify/require"
)

type createDepartmentRepo struct {
	depts      map[int64]models.Department
	nextDeptID int64
	err        error
	gotName    string
	gotParent  *int64
}

func newCreateDepartmentRepo() *createDepartmentRepo {
	parentID := int64(1)
	depts := map[int64]models.Department{
		1: {ID: 1, Name: "Корневой департамент", CreatedAt: time.Now()},
		2: {ID: 2, Name: "Дочерний департамент", ParentID: &parentID, CreatedAt: time.Now()},
	}
	return &createDepartmentRepo{
		depts:      depts,
		nextDeptID: 3,
	}
}

func (r *createDepartmentRepo) CreateDepartment(_ context.Context, name string, parentID *int64) (models.Department, error) {
	r.gotName = name
	r.gotParent = parentID
	if r.err != nil {
		return models.Department{}, r.err
	}

	r.nextDeptID++
	dept := models.Department{
		ID:        r.nextDeptID,
		Name:      name,
		ParentID:  parentID,
		CreatedAt: time.Now(),
	}
	r.depts[dept.ID] = dept
	return dept, nil
}

func (r *createDepartmentRepo) SelectDepartmentById(_ context.Context, id int64) (models.Department, error) {
	dept, ok := r.depts[id]
	if !ok {
		return models.Department{}, models.ErrNotFound
	}
	return dept, nil
}

func (r *createDepartmentRepo) SelectChildrenDepartments(context.Context, int64, int) ([]models.Department, error) {
	return nil, nil
}

func (r *createDepartmentRepo) UpdateDepartment(context.Context, int64, *string, *int64, bool) (models.Department, error) {
	return models.Department{}, nil
}

func Test_CreateDepartment(t *testing.T) {
	parentID := int64(1)

	tests := []struct {
		name     string
		inName   string
		inParent *int64
		want     models.Department
		wantErr  error
	}{
		{
			name:   "успешное_создание",
			inName: "новый департамент",
			want: models.Department{
				ID:   4,
				Name: "новый департамент",
			},
		},
		{
			name:     "успешное_создание_с_parent_id",
			inName:   "новый департамент 2",
			inParent: &parentID,
			want: models.Department{
				ID:       4,
				Name:     "новый департамент 2",
				ParentID: &parentID,
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
			name:   "обрезает_пробелы",
			inName: " новый департамент ",
			want: models.Department{
				ID:   4,
				Name: "новый департамент",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newCreateDepartmentRepo()
			svc := NewService(repo)

			got, err := svc.CreateDepartment(context.Background(), tt.inName, tt.inParent)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, models.Department{}, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.ID, got.ID)
			require.Equal(t, tt.want.Name, got.Name)
			require.Equal(t, tt.want.ParentID, got.ParentID)
			require.False(t, got.CreatedAt.IsZero())
		})
	}
}
