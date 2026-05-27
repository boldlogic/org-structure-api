-- +goose Up
CREATE TABLE org.employees (
    id SERIAL PRIMARY KEY,
    department_id INT NOT NULL REFERENCES org.departments (id) ON DELETE CASCADE,
    full_name VARCHAR(200) NOT NULL,
    position VARCHAR(200) NOT NULL,
    hired_at DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose Down
DROP TABLE IF EXISTS org.employees;

