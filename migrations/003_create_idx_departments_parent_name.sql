-- +goose Up
CREATE UNIQUE INDEX idx_departments_parent_name
    ON org.departments (COALESCE(parent_id, 0), name);
    
-- +goose Down
DROP INDEX IF EXISTS idx_departments_parent_name;
