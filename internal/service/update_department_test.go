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

func Test_UpdateDepartment(t *testing.T) {
	deptID := int64(3)
	parentID := int64(1)
	otherParentID := int64(99)
	newName := "обновлённое"
	created := time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		inName        *string
		inParent      *int64
		parentIDSet   bool
		repoName      *string
		repoParent    *int64
		repoParentSet bool
		repoDep       models.Department
		repoErr       error
		wantDep       models.Department
		wantErr       error
	}{
		{
			name:          "успешное_обновление",
			inName:        &newName,
			repoName:      &newName,
			repoDep: models.Department{
				ID:        deptID,
				Name:      newName,
				CreatedAt: created,
			},
			wantDep: models.Department{
				ID:        deptID,
				Name:      newName,
				CreatedAt: created,
			},
		},
		{
			name:        "parent_id_равен_id",
			inParent:    &deptID,
			parentIDSet: true,
			wantErr:     models.ErrValidation,
		},
		{
			name:          "цикл_в_дереве",
			inParent:      &parentID,
			parentIDSet:   true,
			repoParent:    &parentID,
			repoParentSet: true,
			repoErr:       models.ErrCycle,
			wantErr:       models.ErrConflict,
		},
		{
			name:     "конфликт_названия",
			inName:   &newName,
			repoName: &newName,
			repoErr:  models.ErrConflict,
			wantErr:  models.ErrConflict,
		},
		{
			name:        "не_передано_ни_одно_поле",
			parentIDSet: false,
			wantErr:     models.ErrValidation,
		},
		{
			name:          "parent_id_не_существует",
			inParent:      &otherParentID,
			parentIDSet:   true,
			repoParent:    &otherParentID,
			repoParentSet: true,
			repoErr:       models.ErrParentNotFound,
			wantErr:       models.ErrBusinessValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(testRepo)
			if tt.repoName != nil || tt.repoParent != nil || tt.repoErr != nil {
				repo.On("UpdateDepartment", mock.Anything, deptID, tt.repoName, tt.repoParent, tt.repoParentSet).
					Return(tt.repoDep, tt.repoErr).
					Once()
			}

			svc := NewService(repo)
			got, err := svc.UpdateDepartment(context.Background(), deptID, tt.inName, tt.inParent, tt.parentIDSet)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, models.Department{}, got)
				if tt.repoName == nil && tt.repoParent == nil && tt.repoErr == nil {
					repo.AssertNotCalled(t, "UpdateDepartment")
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
