package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/boldlogic/packages/transport/httputils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createDepartmentDTO(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		want   createDepartmentDTO
		hasErr bool
	}{
		{
			name: "валидный_полный_запрос",
			body: `{"name":"подразделение","parent_id":7}`,
			want: createDepartmentDTO{Name: "подразделение", ParentID: new(int64(7))},
		},
		{
			name: "без_parent_id",
			body: `{"name":"подразделение"}`,
			want: createDepartmentDTO{Name: "подразделение", ParentID: nil},
		},
		{
			name:   "без_name",
			body:   `{"parent_id":7}`,
			hasErr: true,
		},
		{
			name: "валидный_запрос_с_null_parent_id",
			body: `{"name":"подразделение","parent_id":null}`,
			want: createDepartmentDTO{Name: "подразделение", ParentID: nil},
		},
		{
			name:   "name_null",
			body:   `{"name":null}`,
			hasErr: true,
		},
		{
			name: "name_1_символ",
			body: `{"name":"a"}`,
			want: createDepartmentDTO{Name: "a", ParentID: nil},
		},
		{
			name: "name_200_символов",
			body: fmt.Sprintf(`{"name":"%s"}`, strings.Repeat("a", 200)),
			want: createDepartmentDTO{Name: strings.Repeat("a", 200), ParentID: nil},
		},
		{
			name: "parent_id=min",
			body: `{"name":"a","parent_id":1}`,
			want: createDepartmentDTO{Name: "a", ParentID: new(int64(1))},
		},
		{
			name: "parent_id=max",
			body: `{"name":"a","parent_id":2147483647}`,
			want: createDepartmentDTO{Name: "a", ParentID: new(int64(2147483647))},
		},
		{
			name:   "name_выше_ограничения",
			body:   fmt.Sprintf(`{"name":"%s"}`, strings.Repeat("a", 201)),
			hasErr: true,
		},
		{
			name:   "parent_id>max",
			body:   `{"name":"a","parent_id":2147483648}`,
			hasErr: true,
		},
	}
	for _, testCase := range tests {
		tt := testCase

		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			got, err := httputils.DecodeRequest[createDepartmentDTO](req)
			if tt.hasErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)

		})
	}
}

func Test_updateDepartmentDTO(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		check  func(t *testing.T, got updateDepartmentDTO)
		hasErr bool
	}{
		{
			name: "валидный_полный_запрос",
			body: `{"name":"подразделение","parent_id":7}`,
			check: func(t *testing.T, got updateDepartmentDTO) {
				assert.Equal(t, "подразделение", **got.Name)
				assert.Equal(t, int64(7), **got.ParentID)
			},
		},
		{
			name: "name_1_символ",
			body: `{"name":"a"}`,
			check: func(t *testing.T, got updateDepartmentDTO) {
				assert.Equal(t, "a", **got.Name)
			},
		},
		{
			name: "name_200_символов",
			body: fmt.Sprintf(`{"name":"%s"}`, strings.Repeat("a", 200)),
			check: func(t *testing.T, got updateDepartmentDTO) {
				assert.Equal(t, strings.Repeat("a", 200), **got.Name)
			},
		},
		{
			name: "только_parent_id=min",
			body: `{"parent_id":1}`,
			check: func(t *testing.T, got updateDepartmentDTO) {
				assert.Nil(t, *got.Name)
				assert.NotNil(t, got.Name)
				assert.Equal(t, int64(1), **got.ParentID)
			},
		},
		{
			name: "parent_id=max",
			body: `{"parent_id":2147483647}`,
			check: func(t *testing.T, got updateDepartmentDTO) {
				assert.Nil(t, *got.Name)
				assert.NotNil(t, got.Name)
				assert.Equal(t, int64(2147483647), **got.ParentID)
			},
		},
		{
			name: "parent_id_пропущен",
			body: `{"name":"подразделение"}`,
			check: func(t *testing.T, got updateDepartmentDTO) {
				assert.Nil(t, *got.ParentID)
				assert.NotNil(t, got.ParentID)
			},
		},
		{
			name: "null_parent_id",
			body: `{"name":"подразделение","parent_id":null}`,
			check: func(t *testing.T, got updateDepartmentDTO) {
				assert.Equal(t, "подразделение", **got.Name)
				assert.Nil(t, got.ParentID)
			},
		},
		{
			name:   "name_выше_ограничения",
			body:   fmt.Sprintf(`{"name":"%s"}`, strings.Repeat("a", 201)),
			hasErr: true,
		},
		{
			name:   "parent_id>max",
			body:   `{"parent_id":2147483648}`,
			hasErr: true,
		},
	}
	for _, testCase := range tests {
		tt := testCase

		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			got, err := httputils.DecodeRequest[updateDepartmentDTO](req)
			if tt.hasErr {
				require.Error(t, err)
				return
			}
			assert.NoError(t, err)
			tt.check(t, got)

		})
	}
}

func Test_createEmployeeDTO(t *testing.T) {
	date := "2024-06-15"
	tests := []struct {
		name   string
		body   string
		want   createEmployeeDTO
		hasErr bool
	}{
		{
			name: "валидный_полный_запрос",
			body: `{"full_name":"Иванов","position":"инженер","hired_at":"2024-06-15"}`,
			want: createEmployeeDTO{FullName: "Иванов", Position: "инженер", HiredAt: &date},
		},
		{
			name: "без_hired_at",
			body: `{"full_name":"Иванов","position":"инженер"}`,
			want: createEmployeeDTO{FullName: "Иванов", Position: "инженер", HiredAt: nil},
		},
		{
			name:   "без_full_name",
			body:   `{"position":"инженер"}`,
			hasErr: true,
		},
		{
			name:   "без_position",
			body:   `{"full_name":"Иванов"}`,
			hasErr: true,
		},
		{
			name: "валидный_запрос_с_null_hired_at",
			body: `{"full_name":"Иванов","position":"инженер","hired_at":null}`,
			want: createEmployeeDTO{FullName: "Иванов", Position: "инженер", HiredAt: nil},
		},
		{
			name:   "full_name_null",
			body:   `{"full_name":null,"position":"инженер"}`,
			hasErr: true,
		},
		{
			name:   "position_null",
			body:   `{"full_name":"Иванов","position":null}`,
			hasErr: true,
		},
		{
			name: "full_name_1_символ",
			body: `{"full_name":"а","position":"б"}`,
			want: createEmployeeDTO{FullName: "а", Position: "б", HiredAt: nil},
		},
		{
			name: "full_name_200_символов",
			body: fmt.Sprintf(`{"full_name":"%s","position":"б"}`, strings.Repeat("a", 200)),
			want: createEmployeeDTO{FullName: strings.Repeat("a", 200), Position: "б", HiredAt: nil},
		},
		{
			name: "position_200_символов",
			body: fmt.Sprintf(`{"full_name":"а","position":"%s"}`, strings.Repeat("b", 200)),
			want: createEmployeeDTO{FullName: "а", Position: strings.Repeat("b", 200), HiredAt: nil},
		},
		{
			name:   "full_name_выше_ограничения",
			body:   fmt.Sprintf(`{"full_name":"%s","position":"б"}`, strings.Repeat("a", 201)),
			hasErr: true,
		},
		{
			name:   "position_выше_ограничения",
			body:   fmt.Sprintf(`{"full_name":"а","position":"%s"}`, strings.Repeat("b", 201)),
			hasErr: true,
		},
		{
			name:   "hired_at_неверный_формат",
			body:   `{"full_name":"Иванов","position":"инженер","hired_at":"15.06.2024"}`,
			hasErr: true,
		},
		{
			name:   "hired_at_некорректная_дата",
			body:   `{"full_name":"Иванов","position":"инженер","hired_at":"2024-13-40"}`,
			hasErr: true,
		},
	}
	for _, testCase := range tests {
		tt := testCase

		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			got, err := httputils.DecodeRequest[createEmployeeDTO](req)
			if tt.hasErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)

		})
	}
}
