-- +goose Up
CREATE INDEX idx_employees_department_id
ON org.employees (department_id);
    
-- +goose Down
DROP INDEX IF EXISTS org.idx_employees_department_id;
