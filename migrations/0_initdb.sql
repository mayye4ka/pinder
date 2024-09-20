-- +migrate Up
CREATE TABLE users(
    id int NOT NULL AUTO_INCREMENT,
    phone_number varchar(40) NOT NULL,
    pass_hash varchar(512) NOT NULL,
    UNIQUE (phone_number),
    PRIMARY KEY(id)
);
CREATE TABLE profiles(
    user_id int NOT NULL,
    gender varchar(40) NOT NULL,
    age int not null,
    bio text not null,
    photo varchar(80) NOT NULL,
    location_lat double NOT NULL,
    location_lon double NOT NULL,
    location_name varchar(40) NOT NULL,
    UNIQUE (user_id)
);
CREATE TABLE preferences(
    user_id int NOT NULL,
    gender varchar(40) NOT NULL,
    min_age int not null,
    max_age int not null,
    location_lat double NOT NULL,
    location_lon double NOT NULL,
    location_radius_km double NOT NULL,
    UNIQUE (user_id)
);
CREATE TABLE pair_attempts(
    id int NOT NULL AUTO_INCREMENT,
    user1 int NOT NULL,
    user2 int NOT NULL,
    state varchar(40) NOT NULL,
    created_at datetime NOT NULL,
    PRIMARY KEY(id)
);
CREATE TABLE pair_events(
    id int NOT NULL AUTO_INCREMENT,
    paid int NOT NULL,
    event_type varchar(40) NOT NULL,
    created_at datetime NOT NULL,
    PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE users;
DROP TABLE profiles;
DROP TABLE preferences;
DROP TABLE pair_attempts;
DROP TABLE pair_events;
