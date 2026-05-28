package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/boldlogic/org-structure-api/internal/models"
	"github.com/boldlogic/packages/transport/httpserver/response"
	"github.com/boldlogic/packages/transport/httputils"
	"github.com/boldlogic/packages/utils/converters"
	"go.uber.org/zap"
)

type Service interface {
	CreateDepartment(ctx context.Context, name string, parentID *int64) (models.Department, error)
	CreateEmployee(ctx context.Context, emp models.Employee) (models.Employee, error)
	GetDepartment(ctx context.Context, id int64, depth int) (models.Department, []models.Department, error)
	UpdateDepartment(ctx context.Context, id int64, name *string, parentID *int64, parentIDSet bool) (models.Department, error)
	DeleteDepartment(ctx context.Context, id int64) error
}
type Handler struct {
	service Service
	logger  *zap.Logger
}

func NewHandler(svc Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: svc,
		logger:  logger,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("POST /departments/{$}", h.Adapt(h.createDepartment))
	mux.HandleFunc("POST /departments/{id}/employees/{$}", h.Adapt(h.createEmployee))
	mux.HandleFunc("GET /departments/{id}", h.Adapt(h.getDepartment))
	mux.HandleFunc("PATCH /departments/{id}", h.Adapt(h.updateDepartment))
	mux.HandleFunc("DELETE /departments/{id}", h.Adapt(h.deleteDepartment))
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type HandlerFunc func(r *http.Request) (any, string, error)

func (h *Handler) Adapt(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, detail, err := fn(r)
		if err != nil {
			var status int
			switch {
			case errors.Is(err, httputils.ErrUnsupportedMediaType):
				status = http.StatusUnsupportedMediaType
			case errors.Is(err, httputils.ErrRequestEntityTooLarge):
				status = http.StatusRequestEntityTooLarge
			case errors.Is(err, models.ErrValidation) || errors.Is(err, httputils.ErrReadingBody) || errors.Is(err, converters.ErrWrongJSON):
				status = http.StatusBadRequest
			case errors.Is(err, models.ErrBusinessValidation):
				status = http.StatusUnprocessableEntity
			case errors.Is(err, models.ErrConflict):
				status = http.StatusConflict
			case errors.Is(err, models.ErrNotFound):
				status = http.StatusNotFound
			default:
				status = http.StatusInternalServerError
			}
			response.WriteResp(w, status, response.Problem(status, "", detail))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodPost {
			response.WriteResp(w, http.StatusCreated, data)
		} else {
			response.WriteResp(w, http.StatusOK, data)
		}
	}
}
