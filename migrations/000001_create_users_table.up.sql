-- up.sql

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    userid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(225) NOT NULL,
    email VARCHAR(225) NOT NULL UNIQUE,
    password VARCHAR(225) NOT NULL,
    isadmin BOOLEAN DEFAULT false,
    description VARCHAR(255),
    avatar VARCHAR(255)
);

INSERT INTO users (username, email, password) VALUES ('Попкинс', 'popkins@gmail.com', '12345678');
INSERT INTO users (username, email, password) VALUES ('Ванька', 'vanek@gmail.com', '12345678');