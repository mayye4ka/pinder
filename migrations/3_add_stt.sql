-- +migrate Up
CREATE TABLE transcriptions(
    message_id int NOT NULL,
    transcription text NOT NULL,
    PRIMARY KEY(chat_id)
);

-- +migrate Down
DROP TABLE transcriptions;
