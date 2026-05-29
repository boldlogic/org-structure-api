-- +goose Up
CREATE INDEX idx_departments_parent_id
ON org.departments (parent_id);
    
-- +goose Down
DROP INDEX IF EXISTS org.idx_departments_parent_id;
