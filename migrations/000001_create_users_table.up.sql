-- up.sql
CREATE TABLE IF NOT EXISTS users (
    userid SERIAL PRIMARY KEY,
    username VARCHAR(30) NOT NULL,
    email VARCHAR(30) NOT NULL UNIQUE,
    password VARCHAR(30) NOT NULL,
    isadmin BOOLEAN DEFAULT false,
    description VARCHAR(255) DEFAULT NULL,
    avatar VARCHAR(255) DEFAULT NULL
);