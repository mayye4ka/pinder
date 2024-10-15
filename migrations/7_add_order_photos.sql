-- +migrate Up
ALTER TABLE photos ADD COLUMN order INT NOT NULL AFTER photo_key;

-- +migrate Down
ALTER TABLE photos DROP COLUMN order;