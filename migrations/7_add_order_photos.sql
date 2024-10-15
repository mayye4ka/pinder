-- +migrate Up
ALTER TABLE photos ADD COLUMN order_n INT NOT NULL AFTER photo_key;

-- +migrate Down
ALTER TABLE photos DROP COLUMN order_n;