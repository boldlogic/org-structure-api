package server

import (
	"time"

	"github.com/boldlogic/org-structure-api/internal/models"
)

type departmentReqDTO struct {
	Name     string `json:"name" validate:"required,min=1,max=200"`
	ParentID *int64 `json:"parent_id" validate:"omitempty"`
}
type departmentRespDTO struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name,omitempty"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func departmentToDto(dep models.Department) departmentRespDTO {
	return departmentRespDTO{
		ID:        dep.ID,
		Name:      dep.Name,
		ParentID:  dep.ParentID,
		CreatedAt: dep.CreatedAt,
	}
}
