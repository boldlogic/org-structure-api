package server

import (
	"net/http"

	"github.com/boldlogic/packages/transport/httputils"
)

func (h *Handler) createEmployee(r *http.Request) (any, string, error) {
	deptID, err := getDepartmentId(r)
	if err != nil {
		return nil, err.Error(), err
	}

	req, err := httputils.DecodeRequest[createEmployeeDTO](r)
	if err != nil {
		err = mapTransportErrors(err)
		return nil, err.Error(), err
	}

	emp, err := req.toEmployee(deptID)
	if err != nil {
		return nil, err.Error(), err
	}

	got, err := h.service.CreateEmployee(r.Context(), emp)
	if err != nil {
		return nil, err.Error(), err
	}

	return employeeToDto(got), "", nil
}
