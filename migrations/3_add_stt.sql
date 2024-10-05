-- +migrate Up
CREATE TABLE transcriptions(
    message_id int NOT NULL,
    transcription text NOT NULL,
    PRIMARY KEY(message_id)
);

-- +migrate Down
DROP TABLE transcriptions;
