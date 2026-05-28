package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/boldlogic/packages/utils/converters"
	"github.com/boldlogic/packages/utils/dates"
)

type createDepartmentDTO struct {
	Name     string `json:"name" validate:"required,min=1,max=200"`
	ParentID *int64 `json:"parent_id" validate:"omitempty,min=1,max=2147483647"`
}

type updateDepartmentDTO struct {
	Name     **string `json:"name" validate:"omitempty,min=1,max=200"`
	ParentID **int64  `json:"parent_id" validate:"omitempty,min=1,max=2147483647"`
}

type createEmployeeDTO struct {
	FullName string  `json:"full_name" validate:"required,min=1,max=200"`
	Position string  `json:"position" validate:"required,min=1,max=200"`
	HiredAt  *string `json:"hired_at" validate:"omitempty,datetime=2006-01-02"`
}

func (d createEmployeeDTO) toEmployee(deptID int64) (models.Employee, error) {
	var hiredAt *time.Time
	if d.HiredAt != nil {
		var err error
		hiredAt, err = dates.OptionalDatePtr(*d.HiredAt, dates.ISODateFormat)
		if err != nil {
			return models.Employee{}, fmt.Errorf("%w: %w: поле hired_at", models.ErrValidation, err)
		}
	}

	return models.Employee{
		DepartmentID: deptID,
		FullName:     d.FullName,
		Position:     d.Position,
		HiredAt:      hiredAt,
	}, nil
}

func (d *updateDepartmentDTO) UnmarshalJSON(data []byte) error {
	var absentName *string
	var absentParentID *int64
	d.Name = &absentName
	d.ParentID = &absentParentID

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for key := range raw {
		if key != "name" && key != "parent_id" {
			return fmt.Errorf("%w неизвестное поле %q", converters.ErrWrongJSON, key)
		}
	}

	if value, ok := raw["name"]; ok {
		if bytes.Equal(bytes.TrimSpace(value), []byte("null")) {
			d.Name = nil
		} else {
			var name string
			if err := json.Unmarshal(value, &name); err != nil {
				return err
			}
			namePtr := &name
			d.Name = &namePtr
		}
	}

	if value, ok := raw["parent_id"]; ok {
		if bytes.Equal(bytes.TrimSpace(value), []byte("null")) {
			d.ParentID = nil
		} else {
			var parentID int64
			if err := json.Unmarshal(value, &parentID); err != nil {
				return err
			}
			parentIDPtr := &parentID
			d.ParentID = &parentIDPtr
		}
	}

	return nil
}

type departmentRespDTO struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name,omitempty"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type departmentByIdRespDTO struct {
	Department departmentRespDTO   `json:"department"`
	Children   []departmentRespDTO `json:"children,omitempty"`
}

func departmentByIdToDTO(dep models.Department, children []models.Department) departmentByIdRespDTO {
	childs := make([]departmentRespDTO, 0, len(children))

	for _, c := range children {
		childs = append(childs, departmentToDto(c))
	}

	out := departmentByIdRespDTO{
		Department: departmentToDto(dep),
		Children:   childs,
	}
	return out
}

func departmentToDto(dep models.Department) departmentRespDTO {
	return departmentRespDTO{
		ID:        dep.ID,
		Name:      dep.Name,
		ParentID:  dep.ParentID,
		CreatedAt: dep.CreatedAt,
	}
}

type employeeRespDTO struct {
	ID           int64     `json:"id"`
	DepartmentID int64     `json:"department_id"`
	FullName     string    `json:"full_name"`
	Position     string    `json:"position"`
	HiredAt      *string   `json:"hired_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

func employeeToDto(emp models.Employee) employeeRespDTO {
	var hiredAt *string
	if emp.HiredAt != nil {
		s := emp.HiredAt.Format(dates.ISODateFormat)
		hiredAt = &s
	}

	return employeeRespDTO{
		ID:           emp.ID,
		DepartmentID: emp.DepartmentID,
		FullName:     emp.FullName,
		Position:     emp.Position,
		HiredAt:      hiredAt,
		CreatedAt:    emp.CreatedAt,
	}
}
