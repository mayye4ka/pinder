-- +migrate Up
ALTER TABLE profiles ADD COLUMN name varchar(80) NOT NULL AFTER user_id;

-- +migrate Down
ALTER TABLE profiles DROP COLUMN name;
