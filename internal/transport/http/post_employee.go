package server

import (
	"net/http"

	"github.com/boldlogic/packages/transport/httputils"
)

func (h *Handler) createEmployee(r *http.Request) (any, string, error) {
	_, err := getDepartmentId(r)
	if err != nil {
		return nil, err.Error(), err
	}

	req, err := httputils.DecodeRequest[createEmployeeDTO](r)
	if err != nil {
		err = mapTransportErrors(err)
		return nil, err.Error(), err
	}

	return req, "", nil
}
