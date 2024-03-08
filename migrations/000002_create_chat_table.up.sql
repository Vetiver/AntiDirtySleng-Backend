-- up.sql

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS chat (
    chatid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chatName VARCHAR(30) NOT NULL,
    owner UUID,
    FOREIGN KEY (owner) REFERENCES users(userid) ON DELETE SET NULL
);

INSERT INTO chat (chatName, owner) VALUES ('Чат Попкинса и Ванька', (SELECT userid FROM users WHERE username = 'Попкинс'));