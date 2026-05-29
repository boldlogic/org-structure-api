-- +goose Up
DROP INDEX IF EXISTS org.idx_departments_parent_id;

CREATE INDEX idx_departments_parent_id_id
ON org.departments (parent_id, id);

    
-- +goose Down
DROP INDEX IF EXISTS org.idx_departments_parent_id_id;

CREATE INDEX idx_departments_parent_id
ON org.departments (parent_id);