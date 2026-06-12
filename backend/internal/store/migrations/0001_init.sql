-- the +goose Up / +goose Down lines are comments read by goose
-- (migration tool). Up runs when migrating forward, Down is how
-- to undo this migration if we ever roll back

-- +goose Up
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE vehicles (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    make TEXT NOT NULL DEFAULT '',
    model TEXT NOT NULL DEFAULT '',
    year INT,
    mass_kg DOUBLE PRECISION NOT NULL,
    drag_coefficient DOUBLE PRECISION NOT NULL,
    frontal_area_m2 DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE vehicles;
