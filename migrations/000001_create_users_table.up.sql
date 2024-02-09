-- up.sql
CREATE TABLE IF NOT EXISTS users (
    userid SERIAL PRIMARY KEY,
    username VARCHAR(225) NOT NULL,
    email VARCHAR(225) NOT NULL UNIQUE,
    password VARCHAR(225) NOT NULL,
    isadmin BOOLEAN DEFAULT false,
    description VARCHAR(255) DEFAULT NULL,
    avatar VARCHAR(255) DEFAULT NULL
);