package service

import (
	"context"
	"testing"
	"time"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRepo struct {
	depts      map[int64]models.Department
	nextDeptID int64
}

func newtestRepo() *testRepo {
	p := int64(1)
	depts := map[int64]models.Department{
		1: {ID: 1, Name: "test_root", CreatedAt: time.Now()},
		2: {ID: 2, Name: "test_child", ParentID: &p, CreatedAt: time.Now()},
	}
	return &testRepo{
		depts:      depts,
		nextDeptID: 3,
	}
}

func (r *testRepo) CreateDepartment(ctx context.Context, name string, parentID *int64) (models.Department, error) {
	r.nextDeptID++
	d := models.Department{ID: r.nextDeptID, Name: name, ParentID: parentID, CreatedAt: time.Now()}
	r.depts[d.ID] = d
	return d, nil
}

func Test_CreateDepartment(t *testing.T) {
	t.Parallel()

	parent1 := int64(1)

	tests := []struct {
		name     string
		inName   string
		inParent *int64
		want     models.Department
		wantErr  error
	}{
		{
			name:     "создание",
			inName:   "новый департамент",
			inParent: nil,
			want: models.Department{
				ID:   4,
				Name: "новый департамент",
			},
		},
		{
			name:     "создание_с_parent",
			inName:   "новый департамент 2",
			inParent: &parent1,
			want: models.Department{
				ID:       4,
				Name:     "новый департамент 2",
				ParentID: &parent1,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newtestRepo()
			svc := NewService(repo)
			ctx := context.Background()

			got, err := svc.CreateDepartment(ctx, tt.inName, tt.inParent)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, models.Department{}, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Name, got.Name)
				assert.Equal(t, tt.want.ParentID, got.ParentID)
				assert.False(t, got.CreatedAt.IsZero())
			}
		})
	}
}
