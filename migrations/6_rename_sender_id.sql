-- +migrate Up
ALTER TABLE messages RENAME COLUMN sender TO sender_id;

-- +migrate Down
ALTER TABLE messages RENAME COLUMN sender_id TO sender;
