-- +goose Up
ALTER TABLE couriers ADD COLUMN transport_type TEXT NOT NULL DEFAULT 'on_foot';  -- on_foot | scooter | car

-- +goose Down
ALTER TABLE couriers DROP COLUMN transport_type;