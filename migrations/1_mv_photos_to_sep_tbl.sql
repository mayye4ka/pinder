-- +migrate Up
CREATE TABLE photos(
    user_id int NOT NULL AUTO_INCREMENT,
    photo_key varchar(80) NOT NULL,
    PRIMARY KEY(user_id, photo_key)
);
ALTER TABLE profiles DROP COLUMN photo;

-- +migrate Down
DROP TABLE photos;
ALTER TABLE profiles ADD COLUMN photo varchar(80) NOT NULL AFTER bio;
