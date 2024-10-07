-- +migrate Up
ALTER TABLE pair_events RENAME COLUMN paid TO pa_id;

-- +migrate Down
ALTER TABLE pair_events RENAME COLUMN pa_id TO paid;
