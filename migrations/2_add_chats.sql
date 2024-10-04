-- +migrate Up
CREATE TABLE chats(
    id int NOT NULL AUTO_INCREMENT,
    user1 int NOT NULL,
    user2 int NOT NULL,
    PRIMARY KEY(chat_id)
);
CREATE TABLE messages(
    id int NOT NULL AUTO_INCREMENT,
    chat_id int NOT NULL,
    sender int NOT NULL,
    content_type varchar(80) NOT NULL,
    payload text NOT NULL,
    created_at datetime NOT NULL,
    PRIMARY KEY(id),
    KEY(chat_id)
);

-- +migrate Down
DROP TABLE chats;
DROP TABLE messages;
