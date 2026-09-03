-- +goose Up
ALTER TABLE delivery ADD COLUMN status TEXT NOT NULL DEFAULT 'active';  -- active | completed | expired
CREATE UNIQUE INDEX delivery_order_id_active_idx ON delivery (order_id) WHERE status = 'active';

-- +goose Down
DROP INDEX delivery_order_id_active_idx;
ALTER TABLE delivery DROP COLUMN status;